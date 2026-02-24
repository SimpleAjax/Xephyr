package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/internal/repositories"
)

// RealScenarioService implements ScenarioService using database queries
type RealScenarioService struct {
	repos *repositories.Provider
}

// NewRealScenarioService creates a new real scenario service
func NewRealScenarioService(repos *repositories.Provider) ScenarioService {
	return &RealScenarioService{repos: repos}
}

// ListScenarios returns scenarios from database
func (s *RealScenarioService) ListScenarios(ctx context.Context, params dto.ScenarioListQueryParams, orgID string) (*dto.ScenarioListResponse, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, err
	}

	filters := repositories.ScenarioFilters{}
	scenarios, total, err := s.repos.GetScenario().List(ctx, orgUUID, filters, repositories.ListParams{Limit: 100})
	if err != nil {
		return nil, err
	}

	scenarioItems := make([]dto.ScenarioResponse, 0, len(scenarios))
	for _, scenario := range scenarios {
		scenarioItems = append(scenarioItems, s.toScenarioItem(&scenario))
	}

	return &dto.ScenarioListResponse{
		Scenarios: scenarioItems,
		Total:     int(total),
	}, nil
}

// GetScenario returns a single scenario
func (s *RealScenarioService) GetScenario(ctx context.Context, scenarioID string, orgID string) (*dto.ScenarioDetailResponse, error) {
	scenarioUUID, err := uuid.Parse(scenarioID)
	if err != nil {
		return nil, err
	}

	scenario, err := s.repos.GetScenario().GetByID(ctx, scenarioUUID)
	if err != nil {
		return nil, err
	}

	return s.toScenarioDetailResponse(scenario), nil
}

// CreateScenario creates a new scenario
func (s *RealScenarioService) CreateScenario(ctx context.Context, req dto.CreateScenarioRequest, orgID string, userID uuid.UUID) (*dto.ScenarioResponse, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, err
	}

	changesJSON, _ := json.Marshal(req.ProposedChanges)

	scenario := &models.Scenario{
		OrganizationID:  orgUUID,
		Title:           req.Title,
		Description:     req.Description,
		ChangeType:      models.ScenarioChangeType(req.ChangeType),
		Status:          models.ScenarioStatusPending,
		ProposedChanges: models.JSONB{"changes": string(changesJSON)},
		CreatedByID:     userID,
	}

	if err := s.repos.GetScenario().Create(ctx, scenario); err != nil {
		return nil, err
	}

	return &dto.ScenarioResponse{
		ScenarioID:       scenario.ID.String(),
		Title:            scenario.Title,
		ChangeType:       string(scenario.ChangeType),
		Status:           string(scenario.Status),
		CreatedAt:        scenario.CreatedAt,
		SimulationStatus: "pending",
	}, nil
}

// SimulateScenario simulates a scenario
func (s *RealScenarioService) SimulateScenario(ctx context.Context, scenarioID string, req dto.SimulateScenarioRequest, orgID string) (*dto.SimulateScenarioResponse, error) {
	scenarioUUID, err := uuid.Parse(scenarioID)
	if err != nil {
		return nil, err
	}

	scenario, err := s.repos.GetScenario().GetByID(ctx, scenarioUUID)
	if err != nil {
		return nil, err
	}

	// Get impact analysis if exists
	var impact *models.ScenarioImpactAnalysis
	if scenario.ImpactAnalysis != nil {
		impact = scenario.ImpactAnalysis
	}

	affectedProjects := []dto.AffectedProject{}
	affectedTasks := []dto.AffectedTask{}

	// Parse affected projects from impact analysis
	if impact != nil && impact.AffectedProjectIDs != nil {
		if idsJSON, err := json.Marshal(impact.AffectedProjectIDs); err == nil {
			var ids []string
			if err := json.Unmarshal(idsJSON, &ids); err == nil {
				for _, id := range ids {
					affectedProjects = append(affectedProjects, dto.AffectedProject{
						ProjectID:     id,
						Name:          "Project " + id[:8],
						Impact:        "medium",
						DelayDays:     0,
						AffectedTasks: []string{},
					})
				}
			}
		}
	}

	recommendations := []dto.AIRecommendation{}
	if impact != nil && impact.Recommendations != nil {
		if recsJSON, err := json.Marshal(impact.Recommendations); err == nil {
			var recs []string
			if err := json.Unmarshal(recsJSON, &recs); err == nil {
				for i, rec := range recs {
					recommendations = append(recommendations, dto.AIRecommendation{
						Priority:  i + 1,
						Action:    rec,
						Reasoning: "Based on impact analysis",
					})
				}
			}
		}
	}

	// Build timeline comparison
	timeline := dto.TimelineComparison{}
	if impact != nil && impact.TimelineComparison != nil {
		if timelineJSON, err := json.Marshal(impact.TimelineComparison); err == nil {
			var timelineData map[string]interface{}
			json.Unmarshal(timelineJSON, &timelineData)
		}
	}

	// Build cost analysis
	costAnalysis := dto.CostAnalysis{
		TotalCost:  0,
		Confidence: 0.85,
		Breakdown:  []dto.CostBreakdownItem{},
	}
	if impact != nil {
		costAnalysis.TotalCost = impact.CostImpact
	}

	return &dto.SimulateScenarioResponse{
		ScenarioID:       scenarioID,
		SimulationStatus: "completed",
		ImpactAnalysis: dto.ImpactAnalysis{
			AffectedProjects:   affectedProjects,
			AffectedTasks:      affectedTasks,
			TimelineComparison: timeline,
			CostAnalysis:       costAnalysis,
			ResourceImpacts:    []dto.ResourceImpact{},
		},
		AIRecommendations:  recommendations,
		CalculatedAt:       time.Now().UTC(),
		SimulationDuration: "1.2s",
	}, nil
}

// ApplyScenario applies a scenario
func (s *RealScenarioService) ApplyScenario(ctx context.Context, scenarioID string, req dto.ApplyScenarioRequest, orgID string, userID uuid.UUID) (*dto.ApplyScenarioResponse, error) {
	scenarioUUID, err := uuid.Parse(scenarioID)
	if err != nil {
		return nil, err
	}

	if err := s.repos.GetScenario().UpdateStatus(ctx, scenarioUUID, models.ScenarioStatusApplied); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &dto.ApplyScenarioResponse{
		ScenarioID: scenarioID,
		Status:     "applied",
		AppliedAt:  now,
		AppliedBy:  userID,
		Changes: dto.ScenarioChanges{
			TasksReassigned:   []dto.ReassignmentChange{},
			DatesAdjusted:     []dto.DateAdjustment{},
			NotificationsSent: 0,
		},
		FollowUp: dto.ScenarioFollowUp{
			NudgesCreated:         []string{},
			CalendarEventsCreated: false,
		},
	}, nil
}

// RejectScenario rejects a scenario
func (s *RealScenarioService) RejectScenario(ctx context.Context, scenarioID string, orgID string, userID uuid.UUID) (*dto.RejectScenarioResponse, error) {
	scenarioUUID, err := uuid.Parse(scenarioID)
	if err != nil {
		return nil, err
	}

	if err := s.repos.GetScenario().UpdateStatus(ctx, scenarioUUID, models.ScenarioStatusRejected); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &dto.RejectScenarioResponse{
		ScenarioID: scenarioID,
		Status:     string(models.ScenarioStatusRejected),
		RejectedAt: now,
		RejectedBy: userID,
	}, nil
}

// ModifyScenario modifies a scenario
func (s *RealScenarioService) ModifyScenario(ctx context.Context, scenarioID string, req dto.ModifyScenarioRequest, orgID string) (*dto.ScenarioResponse, error) {
	scenarioUUID, err := uuid.Parse(scenarioID)
	if err != nil {
		return nil, err
	}

	scenario, err := s.repos.GetScenario().GetByID(ctx, scenarioUUID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Title != nil && *req.Title != "" {
		scenario.Title = *req.Title
	}
	if req.Description != nil && *req.Description != "" {
		scenario.Description = *req.Description
	}
	if req.ProposedChanges != nil {
		changesJSON, _ := json.Marshal(req.ProposedChanges)
		scenario.ProposedChanges = models.JSONB{"changes": string(changesJSON)}
	}
	scenario.Status = models.ScenarioStatusModified

	if err := s.repos.GetScenario().Update(ctx, scenario); err != nil {
		return nil, err
	}

	return &dto.ScenarioResponse{
		ScenarioID: scenario.ID.String(),
		Title:      scenario.Title,
		ChangeType: string(scenario.ChangeType),
		Status:     string(scenario.Status),
		CreatedAt:  scenario.CreatedAt,
	}, nil
}

// Helper functions

func (s *RealScenarioService) toScenarioItem(scenario *models.Scenario) dto.ScenarioResponse {
	return dto.ScenarioResponse{
		ScenarioID:       scenario.ID.String(),
		Title:            scenario.Title,
		ChangeType:       string(scenario.ChangeType),
		Status:           string(scenario.Status),
		CreatedAt:        scenario.CreatedAt,
		SimulationStatus: string(scenario.Status),
	}
}

func (s *RealScenarioService) toScenarioDetailResponse(scenario *models.Scenario) *dto.ScenarioDetailResponse {
	// Parse proposed changes
	var changes dto.ProposedChanges
	if scenario.ProposedChanges != nil {
		if data, err := json.Marshal(scenario.ProposedChanges); err == nil {
			json.Unmarshal(data, &changes)
		}
	}

	detail := &dto.ScenarioDetailResponse{
		ScenarioResponse: dto.ScenarioResponse{
			ScenarioID:       scenario.ID.String(),
			Title:            scenario.Title,
			ChangeType:       string(scenario.ChangeType),
			Status:           string(scenario.Status),
			CreatedAt:        scenario.CreatedAt,
			SimulationStatus: string(scenario.Status),
			ProposedChanges:  changes,
		},
		Description:       scenario.Description,
		History:           []dto.ScenarioHistoryEntry{},
		AIRecommendations: []dto.AIRecommendation{},
	}

	// Add impact analysis if exists
	if scenario.ImpactAnalysis != nil {
		detail.ImpactAnalysis = &dto.ImpactAnalysis{
			AffectedProjects: []dto.AffectedProject{},
			AffectedTasks:    []dto.AffectedTask{},
			TimelineComparison: dto.TimelineComparison{
				TotalDelayDays: scenario.ImpactAnalysis.DelayHoursTotal / 24,
			},
			CostAnalysis: dto.CostAnalysis{
				TotalCost:  scenario.ImpactAnalysis.CostImpact,
				Confidence: 0.85,
				Breakdown:  []dto.CostBreakdownItem{},
			},
			ResourceImpacts: []dto.ResourceImpact{},
		}
	}

	return detail
}
