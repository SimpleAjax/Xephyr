package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
)

// NudgeService defines the interface for nudge-related operations
type NudgeService interface {
	// ListNudges returns a list of nudges with filters
	ListNudges(ctx context.Context, params dto.NudgeListQueryParams, orgID string, userID uuid.UUID) (*dto.NudgeListResponse, error)

	// GetNudge returns a single nudge by ID
	GetNudge(ctx context.Context, nudgeID string, orgID string) (*dto.NudgeDetailResponse, error)

	// TakeNudgeAction processes an action on a nudge
	TakeNudgeAction(ctx context.Context, nudgeID string, req dto.NudgeActionRequest, orgID string, userID uuid.UUID) (*dto.NudgeActionResponse, error)

	// UpdateNudgeStatus updates the status of a nudge
	UpdateNudgeStatus(ctx context.Context, nudgeID string, status string, orgID string) (*dto.NudgeResponse, error)

	// GenerateNudges triggers manual nudge generation
	GenerateNudges(ctx context.Context, req dto.GenerateNudgesRequest, orgID string) (string, error)

	// GetNudgeStats returns nudge statistics
	GetNudgeStats(ctx context.Context, period string, orgID string) (*dto.NudgeStatsResponse, error)
}

// DummyNudgeService is a placeholder implementation of NudgeService
type DummyNudgeService struct{}

// NewDummyNudgeService creates a new dummy nudge service
func NewDummyNudgeService() NudgeService {
	return &DummyNudgeService{}
}

// ListNudges returns dummy nudges
func (s *DummyNudgeService) ListNudges(ctx context.Context, params dto.NudgeListQueryParams, orgID string, userID uuid.UUID) (*dto.NudgeListResponse, error) {
	return &dto.NudgeListResponse{
		Nudges: []dto.NudgeResponse{
			{
				ID:       "nudge-1",
				Type:     "overload",
				Severity: "high",
				Status:   "unread",
				Title:    "Emma Wilson is overallocated at 125%",
				Description: "Emma has 75 hours assigned this week across 3 projects",
				AIExplanation: "Based on current assignments, Emma is scheduled for 75 hours this week...",
				SuggestedAction: dto.SuggestedAction{
					Type:                "reassign",
					Description:         "Reassign Marketing Website consultation to Rachel",
					TargetTaskID:        strPtr("task-web-3"),
					SuggestedAssigneeID: strPtr("user-rachel"),
				},
				RelatedEntities: dto.RelatedEntities{
					ProjectID: strPtr("proj-mobile"),
					TaskID:    strPtr("task-fit-1"),
					PersonID:  strPtr("user-emma"),
				},
				Metrics: dto.NudgeMetrics{
					AllocationPercentage: 125,
					AssignedTasks:        4,
					TotalHours:           75,
				},
				CriticalityScore: 90,
				CreatedAt:        time.Now().UTC().Add(-24 * time.Hour),
				ExpiresAt:        timePtr(time.Now().UTC().Add(7 * 24 * time.Hour)),
			},
		},
		Summary: dto.NudgeSummary{
			Total:  12,
			Unread: 5,
			BySeverity: map[string]int{
				"high":   3,
				"medium": 5,
				"low":    4,
			},
			ByType: map[string]int{
				"overload":   2,
				"delay_risk": 3,
				"unassigned": 2,
			},
		},
	}, nil
}

// GetNudge returns a dummy nudge
func (s *DummyNudgeService) GetNudge(ctx context.Context, nudgeID string, orgID string) (*dto.NudgeDetailResponse, error) {
	return &dto.NudgeDetailResponse{
		NudgeResponse: dto.NudgeResponse{
			ID:       nudgeID,
			Type:     "overload",
			Severity: "high",
			Status:   "unread",
			Title:    "Sample Nudge",
			Description: "This is a sample nudge",
			CriticalityScore: 90,
			CreatedAt:        time.Now().UTC(),
		},
		History: []dto.NudgeAction{
			{
				Action:    "created",
				Timestamp: time.Now().UTC().Add(-24 * time.Hour),
				UserID:    uuid.New(),
			},
		},
	}, nil
}

// TakeNudgeAction processes a dummy action
func (s *DummyNudgeService) TakeNudgeAction(ctx context.Context, nudgeID string, req dto.NudgeActionRequest, orgID string, userID uuid.UUID) (*dto.NudgeActionResponse, error) {
	return &dto.NudgeActionResponse{
		NudgeID:     nudgeID,
		ActionTaken: req.ActionType,
		Result: dto.NudgeActionResult{
			TaskReassigned: true,
			FromUserID:     "user-emma",
			ToUserID:       "user-rachel",
			TaskID:         "task-web-3",
		},
		NudgeStatus:   "acted",
		FollowUpNudges: []string{},
		CompletedAt:   time.Now().UTC(),
	}, nil
}

// UpdateNudgeStatus updates dummy nudge status
func (s *DummyNudgeService) UpdateNudgeStatus(ctx context.Context, nudgeID string, status string, orgID string) (*dto.NudgeResponse, error) {
	return &dto.NudgeResponse{
		ID:     nudgeID,
		Type:   "overload",
		Status: status,
		Title:  "Sample Nudge",
	}, nil
}

// GenerateNudges triggers dummy nudge generation
func (s *DummyNudgeService) GenerateNudges(ctx context.Context, req dto.GenerateNudgesRequest, orgID string) (string, error) {
	return "job_nudge_gen_abc123", nil
}

// GetNudgeStats returns dummy stats
func (s *DummyNudgeService) GetNudgeStats(ctx context.Context, period string, orgID string) (*dto.NudgeStatsResponse, error) {
	return &dto.NudgeStatsResponse{
		Period:          period,
		Generated:       45,
		Acted:           23,
		Dismissed:       12,
		Expired:         5,
		ActionRate:      0.51,
		AvgTimeToAction: "4.2h",
		ByType: map[string]dto.NudgeTypeStats{
			"overload":   {Generated: 8, Acted: 6},
			"delay_risk": {Generated: 12, Acted: 7},
		},
	}, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}
