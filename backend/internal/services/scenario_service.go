package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
)

// ScenarioService defines the interface for scenario-related operations
type ScenarioService interface {
	// CreateScenario creates a new scenario
	CreateScenario(ctx context.Context, req dto.CreateScenarioRequest, orgID string, createdBy uuid.UUID) (*dto.ScenarioResponse, error)

	// ListScenarios returns a list of scenarios
	ListScenarios(ctx context.Context, params dto.ScenarioListQueryParams, orgID string) (*dto.ScenarioListResponse, error)

	// GetScenario returns a single scenario by ID
	GetScenario(ctx context.Context, scenarioID string, orgID string) (*dto.ScenarioDetailResponse, error)

	// SimulateScenario runs a simulation for a scenario
	SimulateScenario(ctx context.Context, scenarioID string, req dto.SimulateScenarioRequest, orgID string) (*dto.SimulateScenarioResponse, error)

	// ApplyScenario applies a scenario
	ApplyScenario(ctx context.Context, scenarioID string, req dto.ApplyScenarioRequest, orgID string, appliedBy uuid.UUID) (*dto.ApplyScenarioResponse, error)

	// RejectScenario rejects a scenario
	RejectScenario(ctx context.Context, scenarioID string, orgID string, rejectedBy uuid.UUID) (*dto.RejectScenarioResponse, error)

	// ModifyScenario modifies a scenario
	ModifyScenario(ctx context.Context, scenarioID string, req dto.ModifyScenarioRequest, orgID string) (*dto.ScenarioResponse, error)
}

// DummyScenarioService is a placeholder implementation of ScenarioService
type DummyScenarioService struct{}

// NewDummyScenarioService creates a new dummy scenario service
func NewDummyScenarioService() ScenarioService {
	return &DummyScenarioService{}
}

// CreateScenario creates a dummy scenario
func (s *DummyScenarioService) CreateScenario(ctx context.Context, req dto.CreateScenarioRequest, orgID string, createdBy uuid.UUID) (*dto.ScenarioResponse, error) {
	return &dto.ScenarioResponse{
		ScenarioID:     "scenario-123",
		Title:          req.Title,
		ChangeType:     req.ChangeType,
		Status:         "draft",
		ProposedChanges: req.ProposedChanges,
		CreatedAt:      time.Now().UTC(),
		SimulationStatus: "pending",
	}, nil
}

// ListScenarios returns dummy scenarios
func (s *DummyScenarioService) ListScenarios(ctx context.Context, params dto.ScenarioListQueryParams, orgID string) (*dto.ScenarioListResponse, error) {
	return &dto.ScenarioListResponse{
		Scenarios: []dto.ScenarioResponse{
			{
				ScenarioID:     "scenario-123",
				Title:          "Emma takes 1-week vacation",
				ChangeType:     "employee_leave",
				Status:         "draft",
				CreatedAt:      time.Now().UTC().Add(-24 * time.Hour),
				SimulationStatus: "completed",
			},
		},
		Total: 1,
	}, nil
}

// GetScenario returns a dummy scenario
func (s *DummyScenarioService) GetScenario(ctx context.Context, scenarioID string, orgID string) (*dto.ScenarioDetailResponse, error) {
	return &dto.ScenarioDetailResponse{
		ScenarioResponse: dto.ScenarioResponse{
			ScenarioID:     scenarioID,
			Title:          "Emma takes 1-week vacation",
			ChangeType:     "employee_leave",
			Status:         "draft",
			CreatedAt:      time.Now().UTC(),
			SimulationStatus: "completed",
		},
		Description: "Simulate impact of Emma taking vacation next week",
		ImpactAnalysis: &dto.ImpactAnalysis{
			AffectedProjects: []dto.AffectedProject{
				{
					ProjectID:     "proj-mobile",
					Name:          "Fitness App Mobile Launch",
					Impact:        "high",
					DelayDays:     5,
					AffectedTasks: []string{"task-fit-1"},
				},
			},
			AffectedTasks: []dto.AffectedTask{
				{
					TaskID:          "task-fit-1",
					Title:           "Mobile App UI Design",
					OriginalDueDate: timePtr(time.Now().UTC().Add(7 * 24 * time.Hour)),
					NewDueDate:      timePtr(time.Now().UTC().Add(14 * 24 * time.Hour)),
					DelayDays:       5,
					Reason:          "Primary assignee on leave",
				},
			},
			TimelineComparison: dto.TimelineComparison{
				TotalDelayDays: 7,
			},
			CostAnalysis: dto.CostAnalysis{
				TotalCost:  3200,
				Confidence: 0.85,
			},
		},
		History: []dto.ScenarioHistoryEntry{
			{
				Action:    "created",
				Timestamp: time.Now().UTC().Add(-24 * time.Hour),
				UserID:    uuid.New(),
			},
		},
	}, nil
}

// SimulateScenario returns dummy simulation
func (s *DummyScenarioService) SimulateScenario(ctx context.Context, scenarioID string, req dto.SimulateScenarioRequest, orgID string) (*dto.SimulateScenarioResponse, error) {
	return &dto.SimulateScenarioResponse{
		ScenarioID:       scenarioID,
		SimulationStatus: "completed",
		ImpactAnalysis: dto.ImpactAnalysis{
			AffectedProjects: []dto.AffectedProject{
				{
					ProjectID:     "proj-mobile",
					Name:          "Fitness App Mobile Launch",
					Impact:        "high",
					DelayDays:     5,
					AffectedTasks: []string{"task-fit-1"},
				},
			},
			TimelineComparison: dto.TimelineComparison{
				OriginalEndDate: timePtr(time.Now().UTC().Add(60 * 24 * time.Hour)),
				NewEndDate:      timePtr(time.Now().UTC().Add(67 * 24 * time.Hour)),
				TotalDelayDays:  7,
			},
			CostAnalysis: dto.CostAnalysis{
				TotalCost: 3200,
				Breakdown: []dto.CostBreakdownItem{
					{Category: "Delay Cost", Amount: 2400},
					{Category: "Context Switch", Amount: 400},
					{Category: "Coverage Cost", Amount: 400},
				},
				Confidence: 0.85,
			},
		},
		AIRecommendations: []dto.AIRecommendation{
			{
				Priority:        1,
				Action:          "Reassign Marketing Website consultation to Rachel immediately",
				Reasoning:       "This task has flexibility in timeline and Rachel has relevant skills",
				EstimatedImpact: "Reduces Fitness App delay from 5 to 3 days",
			},
		},
		CalculatedAt:       time.Now().UTC(),
		SimulationDuration: "2.3s",
	}, nil
}

// ApplyScenario applies dummy scenario
func (s *DummyScenarioService) ApplyScenario(ctx context.Context, scenarioID string, req dto.ApplyScenarioRequest, orgID string, appliedBy uuid.UUID) (*dto.ApplyScenarioResponse, error) {
	return &dto.ApplyScenarioResponse{
		ScenarioID: scenarioID,
		Status:     "applied",
		AppliedAt:  time.Now().UTC(),
		AppliedBy:  appliedBy,
		Changes: dto.ScenarioChanges{
			TasksReassigned: []dto.ReassignmentChange{
				{
					TaskID: "task-web-3",
					From:   "user-emma",
					To:     "user-rachel",
				},
			},
			DatesAdjusted: []dto.DateAdjustment{
				{
					TaskID:     "task-fit-1",
					NewDueDate: timePtr(time.Now().UTC().Add(10 * 24 * time.Hour)),
				},
			},
			NotificationsSent: 5,
		},
		FollowUp: dto.ScenarioFollowUp{
			NudgesCreated:       []string{"nudge-reallocation-1"},
			CalendarEventsCreated: true,
		},
	}, nil
}

// RejectScenario rejects dummy scenario
func (s *DummyScenarioService) RejectScenario(ctx context.Context, scenarioID string, orgID string, rejectedBy uuid.UUID) (*dto.RejectScenarioResponse, error) {
	return &dto.RejectScenarioResponse{
		ScenarioID: scenarioID,
		Status:     "rejected",
		RejectedAt: time.Now().UTC(),
		RejectedBy: rejectedBy,
	}, nil
}

// ModifyScenario modifies dummy scenario
func (s *DummyScenarioService) ModifyScenario(ctx context.Context, scenarioID string, req dto.ModifyScenarioRequest, orgID string) (*dto.ScenarioResponse, error) {
	return &dto.ScenarioResponse{
		ScenarioID: scenarioID,
		Title:      "Modified Scenario",
		ChangeType: "employee_leave",
		Status:     "modified",
		CreatedAt:  time.Now().UTC(),
	}, nil
}
