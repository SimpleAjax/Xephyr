package services

import (
	"context"
	"time"

	"github.com/SimpleAjax/Xephyr/internal/dto"
)

// DependencyService defines the interface for dependency-related operations
type DependencyService interface {
	// GetTaskDependencies returns all dependencies for a task
	GetTaskDependencies(ctx context.Context, taskID string, includeIndirect bool, orgID string) (*dto.TaskDependenciesResponse, error)

	// CreateDependency creates a new dependency
	CreateDependency(ctx context.Context, req dto.CreateDependencyRequest, orgID string) (*dto.CreateDependencyResponse, error)

	// DeleteDependency removes a dependency
	DeleteDependency(ctx context.Context, dependencyID string, orgID string) error

	// GetCriticalPath returns the critical path for a project
	GetCriticalPath(ctx context.Context, projectID string, orgID string) (*dto.CriticalPathResponse, error)

	// ValidateDependency validates a potential dependency
	ValidateDependency(ctx context.Context, req dto.ValidateDependencyRequest, orgID string) (*dto.ValidateDependencyResponse, error)

	// GetDependencyGraph returns the dependency graph for a project
	GetDependencyGraph(ctx context.Context, projectID string, orgID string) (*dto.DependencyGraphResponse, error)
}

// DummyDependencyService is a placeholder implementation of DependencyService
type DummyDependencyService struct{}

// NewDummyDependencyService creates a new dummy dependency service
func NewDummyDependencyService() DependencyService {
	return &DummyDependencyService{}
}

// GetTaskDependencies returns dummy dependencies
func (s *DummyDependencyService) GetTaskDependencies(ctx context.Context, taskID string, includeIndirect bool, orgID string) (*dto.TaskDependenciesResponse, error) {
	resp := &dto.TaskDependenciesResponse{
		TaskID: taskID,
		Dependencies: dto.DependencySection{
			Direct: []dto.DependencyInfo{
				{
					DependencyID:    "dep-1",
					DependsOnTaskID: "task-ec-2",
					DependencyType:  "finish_to_start",
					LagHours:        0,
					Status:          "in_progress",
					IsBlocking:      true,
				},
			},
			Indirect: []dto.IndirectDependency{},
		},
		Dependents: dto.DependentsSection{
			Direct: []dto.DependentInfo{
				{
					TaskID:         "task-ec-5",
					DependencyType: "finish_to_start",
					IsBlocked:      true,
				},
			},
			Indirect: []dto.IndirectDependency{},
		},
		ChainAnalysis: dto.ChainAnalysis{
			LongestChain:         3,
			CriticalPathPosition: "on_path",
			FloatHours:           0,
		},
	}

	if includeIndirect {
		resp.Dependencies.Indirect = []dto.IndirectDependency{
			{
				Path:  []string{"task-ec-4", "task-ec-2", "task-ec-1"},
				Depth: 2,
			},
		}
		resp.Dependents.Indirect = []dto.IndirectDependency{
			{
				Path:  []string{"task-ec-4", "task-ec-5", "task-ec-6"},
				Depth: 2,
			},
		}
	}

	return resp, nil
}

// CreateDependency creates a dummy dependency
func (s *DummyDependencyService) CreateDependency(ctx context.Context, req dto.CreateDependencyRequest, orgID string) (*dto.CreateDependencyResponse, error) {
	// Check for circular dependency scenario (task-ec-1 depending on task-ec-4)
	if req.TaskID == "task-ec-1" && req.DependsOnTaskID == "task-ec-4" {
		return nil, &CircularDependencyError{
			Message: "Circular dependency detected",
			Cycle:   []string{"task-ec-1", "task-ec-4", "task-ec-2", "task-ec-1"},
		}
	}
	
	return &dto.CreateDependencyResponse{
		DependencyID:    "dep-123",
		TaskID:          req.TaskID,
		DependsOnTaskID: req.DependsOnTaskID,
		DependencyType:  req.DependencyType,
		LagHours:        req.LagHours,
		CreatedAt:       time.Now().UTC(),
		Validation: dto.DependencyValidation{
			Valid:            true,
			WouldCreateCycle: false,
		},
		Impact: dto.DependencyImpact{
			CriticalPathChanged: true,
			AffectedTasks:       []string{req.TaskID, "task-ec-5"},
			NewProjectEndDate:   timePtr(time.Now().UTC().Add(30 * 24 * time.Hour)),
		},
	}, nil
}

// CircularDependencyError represents a circular dependency error
type CircularDependencyError struct {
	Message string
	Cycle   []string
}

func (e *CircularDependencyError) Error() string {
	return e.Message
}

// DeleteDependency is a no-op for dummy service
func (s *DummyDependencyService) DeleteDependency(ctx context.Context, dependencyID string, orgID string) error {
	return nil
}

// GetCriticalPath returns dummy critical path
func (s *DummyDependencyService) GetCriticalPath(ctx context.Context, projectID string, orgID string) (*dto.CriticalPathResponse, error) {
	return &dto.CriticalPathResponse{
		ProjectID: projectID,
		CriticalPath: dto.CriticalPathInfo{
			TaskIDs: []string{"task-ec-1", "task-ec-2", "task-ec-4"},
			Tasks: []dto.CriticalPathTask{
				{
					TaskID:         "task-ec-1",
					Title:          "Design System Architecture",
					EstimatedHours: 40,
					EarliestStart:  timePtr(time.Now().UTC().Add(-30 * 24 * time.Hour)),
					EarliestFinish: timePtr(time.Now().UTC().Add(-20 * 24 * time.Hour)),
					LatestStart:    timePtr(time.Now().UTC().Add(-30 * 24 * time.Hour)),
					LatestFinish:   timePtr(time.Now().UTC().Add(-20 * 24 * time.Hour)),
					FloatHours:     0,
				},
				{
					TaskID:         "task-ec-2",
					Title:          "Backend API Development",
					EstimatedHours: 80,
					FloatHours:     0,
				},
				{
					TaskID:         "task-ec-4",
					Title:          "Checkout Flow Implementation",
					EstimatedHours: 50,
					FloatHours:     0,
				},
			},
			TotalDuration: 170,
		},
		NonCriticalTasks: []dto.NonCriticalTask{
			{
				TaskID:     "task-ec-5",
				Title:      "Admin Dashboard",
				FloatHours: 40,
			},
		},
		ProjectDuration: 170,
		CalculatedAt:    time.Now().UTC(),
	}, nil
}

// ValidateDependency returns dummy validation
func (s *DummyDependencyService) ValidateDependency(ctx context.Context, req dto.ValidateDependencyRequest, orgID string) (*dto.ValidateDependencyResponse, error) {
	return &dto.ValidateDependencyResponse{
		Valid:            true,
		WouldCreateCycle: false,
		Warnings: []dto.DependencyWarning{
			{
				Type:    "date_constraint",
				Message: "Task due date may need adjustment based on dependency",
			},
		},
		Impact: dto.ValidationImpact{
			EstimatedDelay: 0,
			AffectedTasks:  1,
		},
	}, nil
}

// GetDependencyGraph returns dummy graph
func (s *DummyDependencyService) GetDependencyGraph(ctx context.Context, projectID string, orgID string) (*dto.DependencyGraphResponse, error) {
	return &dto.DependencyGraphResponse{
		ProjectID: projectID,
		Nodes: []dto.DependencyGraphNode{
			{ID: "task-1", Title: "Task 1", Status: "done", X: 0, Y: 0},
			{ID: "task-2", Title: "Task 2", Status: "in_progress", X: 100, Y: 0},
			{ID: "task-3", Title: "Task 3", Status: "backlog", X: 200, Y: 0},
		},
		Edges: []dto.DependencyGraphEdge{
			{ID: "dep-1", Source: "task-1", Target: "task-2", Type: "finish_to_start", LagHours: 0},
			{ID: "dep-2", Source: "task-2", Target: "task-3", Type: "finish_to_start", LagHours: 8},
		},
	}, nil
}
