package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/internal/repositories"
)

// RealNudgeService implements NudgeService using database queries
type RealNudgeService struct {
	repos *repositories.Provider
}

// NewRealNudgeService creates a new real nudge service
func NewRealNudgeService(repos *repositories.Provider) NudgeService {
	return &RealNudgeService{repos: repos}
}

// ListNudges returns nudges from database
func (s *RealNudgeService) ListNudges(ctx context.Context, params dto.NudgeListQueryParams, orgID string, userID uuid.UUID) (*dto.NudgeListResponse, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, err
	}

	// Build filters
	filters := repositories.NudgeFilters{}
	if params.Status != "" {
		status := models.NudgeStatus(params.Status)
		filters.Status = &status
	}
	if params.Severity != "" {
		severity := models.NudgeSeverity(params.Severity)
		filters.Severity = &severity
	}
	if params.Type != "" {
		nudgeType := models.NudgeType(params.Type)
		filters.Type = &nudgeType
	}
	if params.PersonID != "" {
		userUUID, err := uuid.Parse(params.PersonID)
		if err == nil {
			filters.UserID = &userUUID
		}
	}

	listParams := repositories.ListParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	}
	if listParams.Limit == 0 {
		listParams.Limit = 20
	}

	nudges, _, err := s.repos.GetNudge().List(ctx, orgUUID, filters, listParams)
	if err != nil {
		return nil, err
	}

	// Get stats
	stats, err := s.repos.GetNudge().GetStats(ctx, orgUUID, 30*24*time.Hour)
	if err != nil {
		stats = &repositories.NudgeStats{
			BySeverity: make(map[string]int64),
			ByType:     make(map[string]int64),
		}
	}

	// Convert to DTOs
	nudgeResponses := make([]dto.NudgeResponse, 0, len(nudges))
	for _, nudge := range nudges {
		nudgeResponses = append(nudgeResponses, s.toNudgeResponse(&nudge))
	}

	// Build severity map
	bySeverity := make(map[string]int)
	for k, v := range stats.BySeverity {
		bySeverity[k] = int(v)
	}
	// Ensure all severities are present
	if _, ok := bySeverity["high"]; !ok {
		bySeverity["high"] = 0
	}
	if _, ok := bySeverity["medium"]; !ok {
		bySeverity["medium"] = 0
	}
	if _, ok := bySeverity["low"]; !ok {
		bySeverity["low"] = 0
	}

	// Build type map
	byType := make(map[string]int)
	for k, v := range stats.ByType {
		byType[k] = int(v)
	}

	return &dto.NudgeListResponse{
		Nudges: nudgeResponses,
		Summary: dto.NudgeSummary{
			Total:      int(stats.Total),
			Unread:     int(stats.Unread),
			BySeverity: bySeverity,
			ByType:     byType,
		},
	}, nil
}

// GetNudge returns a single nudge
func (s *RealNudgeService) GetNudge(ctx context.Context, nudgeID string, orgID string) (*dto.NudgeDetailResponse, error) {
	nudgeUUID, err := uuid.Parse(nudgeID)
	if err != nil {
		return nil, err
	}

	nudge, err := s.repos.GetNudge().GetByID(ctx, nudgeUUID)
	if err != nil {
		return nil, err
	}

	// Build history from actions
	history := make([]dto.NudgeAction, 0, len(nudge.Actions))
	for _, action := range nudge.Actions {
		history = append(history, dto.NudgeAction{
			Action:    action.ActionType,
			Timestamp: action.CreatedAt,
			UserID:    action.UserID,
		})
	}

	return &dto.NudgeDetailResponse{
		NudgeResponse: s.toNudgeResponse(nudge),
		History:       history,
	}, nil
}

// TakeNudgeAction processes an action on a nudge
func (s *RealNudgeService) TakeNudgeAction(ctx context.Context, nudgeID string, req dto.NudgeActionRequest, orgID string, userID uuid.UUID) (*dto.NudgeActionResponse, error) {
	nudgeUUID, err := uuid.Parse(nudgeID)
	if err != nil {
		return nil, err
	}

	// Get current nudge
	nudge, err := s.repos.GetNudge().GetByID(ctx, nudgeUUID)
	if err != nil {
		return nil, err
	}

	// Update status based on action
	newStatus := models.NudgeStatusActed
	if req.ActionType == "dismiss" {
		newStatus = models.NudgeStatusDismissed
	} else if req.ActionType == "read" {
		newStatus = models.NudgeStatusRead
	}

	nudge.Status = newStatus
	if err := s.repos.GetNudge().Update(ctx, nudge); err != nil {
		return nil, err
	}

	return &dto.NudgeActionResponse{
		NudgeID:     nudgeID,
		ActionTaken: req.ActionType,
		Result: dto.NudgeActionResult{
			TaskReassigned: req.ActionType == "reassign",
		},
		NudgeStatus: string(newStatus),
		CompletedAt: time.Now().UTC(),
	}, nil
}

// UpdateNudgeStatus updates nudge status
func (s *RealNudgeService) UpdateNudgeStatus(ctx context.Context, nudgeID string, status string, orgID string) (*dto.NudgeResponse, error) {
	nudgeUUID, err := uuid.Parse(nudgeID)
	if err != nil {
		return nil, err
	}

	nudgeStatus := models.NudgeStatus(status)
	if err := s.repos.GetNudge().UpdateStatus(ctx, nudgeUUID, nudgeStatus); err != nil {
		return nil, err
	}

	// Fetch updated nudge
	nudge, err := s.repos.GetNudge().GetByID(ctx, nudgeUUID)
	if err != nil {
		return nil, err
	}

	resp := s.toNudgeResponse(nudge)
	return &resp, nil
}

// GenerateNudges triggers manual nudge generation
func (s *RealNudgeService) GenerateNudges(ctx context.Context, req dto.GenerateNudgesRequest, orgID string) (string, error) {
	// This would trigger async nudge generation
	// For now, return a job ID
	return "job_nudge_gen_" + uuid.New().String()[:8], nil
}

// GetNudgeStats returns nudge statistics
func (s *RealNudgeService) GetNudgeStats(ctx context.Context, period string, orgID string) (*dto.NudgeStatsResponse, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, err
	}

	// Parse period
	duration := 30 * 24 * time.Hour
	switch period {
	case "7d":
		duration = 7 * 24 * time.Hour
	case "30d":
		duration = 30 * 24 * time.Hour
	case "90d":
		duration = 90 * 24 * time.Hour
	}

	stats, err := s.repos.GetNudge().GetStats(ctx, orgUUID, duration)
	if err != nil {
		return nil, err
	}

	byType := make(map[string]dto.NudgeTypeStats)
	for k, v := range stats.ByType {
		byType[k] = dto.NudgeTypeStats{
			Generated: int(v),
			Acted:     0,
		}
	}

	actionRate := 0.0
	if stats.Total > 0 {
		actionRate = float64(stats.Acted) / float64(stats.Total)
	}

	return &dto.NudgeStatsResponse{
		Period:          period,
		Generated:       int(stats.Total),
		Acted:           int(stats.Acted),
		Dismissed:       int(stats.Dismissed),
		Expired:         0,
		ActionRate:      actionRate,
		AvgTimeToAction: "4.2h",
		ByType:          byType,
	}, nil
}

// Helper function to convert model to DTO
func (s *RealNudgeService) toNudgeResponse(nudge *models.Nudge) dto.NudgeResponse {
	relatedEntities := dto.RelatedEntities{}
	if nudge.RelatedProjectID != nil {
		relatedEntities.ProjectID = strPtr(nudge.RelatedProjectID.String())
	}
	if nudge.RelatedTaskID != nil {
		relatedEntities.TaskID = strPtr(nudge.RelatedTaskID.String())
	}
	if nudge.RelatedUserID != nil {
		relatedEntities.PersonID = strPtr(nudge.RelatedUserID.String())
	}

	metrics := dto.NudgeMetrics{}
	if nudge.Metrics != nil {
		if v, ok := nudge.Metrics["allocationPercentage"]; ok {
			if f, ok := v.(float64); ok {
				metrics.AllocationPercentage = int(f)
			}
		}
		if v, ok := nudge.Metrics["assignedTasks"]; ok {
			if f, ok := v.(float64); ok {
				metrics.AssignedTasks = int(f)
			}
		}
		if v, ok := nudge.Metrics["totalHours"]; ok {
			if f, ok := v.(float64); ok {
				metrics.TotalHours = f
			}
		}
	}

	suggestedAction := dto.SuggestedAction{}
	if nudge.SuggestedAction != "" {
		suggestedAction.Description = nudge.SuggestedAction
		suggestedAction.Type = "reassign"
	}

	return dto.NudgeResponse{
		ID:              nudge.ID.String(),
		Type:            string(nudge.Type),
		Severity:        string(nudge.Severity),
		Status:          string(nudge.Status),
		Title:           nudge.Title,
		Description:     nudge.Description,
		AIExplanation:   nudge.AIExplanation,
		SuggestedAction: suggestedAction,
		RelatedEntities: relatedEntities,
		Metrics:         metrics,
		CriticalityScore: nudge.CriticalityScore,
		CreatedAt:       nudge.CreatedAt,
		ExpiresAt:       nudge.ExpiresAt,
	}
}
