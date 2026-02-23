package services

import (
	"context"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
)

// PriorityService defines the interface for priority-related operations
type PriorityService interface {
	// GetTaskPriority returns the priority for a single task
	GetTaskPriority(ctx context.Context, taskID string, orgID string) (*dto.TaskPriorityResponse, error)

	// GetBulkTaskPriorities returns priorities for multiple tasks
	GetBulkTaskPriorities(ctx context.Context, req dto.BulkPriorityRequest, orgID string) (*dto.BulkPriorityResponse, error)

	// GetProjectTaskRanking returns the task ranking for a project
	GetProjectTaskRanking(ctx context.Context, projectID string, params dto.ProjectRankingQueryParams, orgID string) (*dto.ProjectTaskRankingResponse, error)

	// RecalculatePriorities triggers priority recalculation
	RecalculatePriorities(ctx context.Context, req dto.RecalculatePriorityRequest, orgID string) (*dto.RecalculatePrioritySyncResponse, *dto.RecalculatePriorityAsyncResponse, error)

	// GetTaskRank returns just the rank for a task
	GetTaskRank(ctx context.Context, taskID string, orgID string) (int, error)
}

// DummyPriorityService is a placeholder implementation of PriorityService
type DummyPriorityService struct{}

// NewDummyPriorityService creates a new dummy priority service
func NewDummyPriorityService() PriorityService {
	return &DummyPriorityService{}
}

// GetTaskPriority returns a dummy task priority
func (s *DummyPriorityService) GetTaskPriority(ctx context.Context, taskID string, orgID string) (*dto.TaskPriorityResponse, error) {
	return &dto.TaskPriorityResponse{
		TaskID:                 taskID,
		UniversalPriorityScore: 87,
		Rank:                   3,
		Breakdown: dto.PriorityBreakdown{
			ProjectPriority:    24,
			BusinessValue:      25,
			DeadlineUrgency:    20,
			CriticalPathWeight: 15,
			DependencyImpact:   3,
		},
		Factors: dto.PriorityFactors{
			ProjectPriority:   95,
			BusinessValue:     100,
			DeadlineUrgency:   80,
			IsOnCriticalPath:  true,
			BlockedTasksCount: 2,
			DaysUntilDue:      42,
		},
	}, nil
}

// GetBulkTaskPriorities returns dummy bulk priorities
func (s *DummyPriorityService) GetBulkTaskPriorities(ctx context.Context, req dto.BulkPriorityRequest, orgID string) (*dto.BulkPriorityResponse, error) {
	priorities := make([]dto.TaskPrioritySummary, len(req.TaskIds))
	for i, taskID := range req.TaskIds {
		priorities[i] = dto.TaskPrioritySummary{
			TaskID:                 taskID,
			UniversalPriorityScore: 50 + i*10,
			Rank:                   i + 1,
		}
	}
	return &dto.BulkPriorityResponse{
		Priorities:  priorities,
		SortedOrder: req.TaskIds,
	}, nil
}

// GetProjectTaskRanking returns dummy project ranking
func (s *DummyPriorityService) GetProjectTaskRanking(ctx context.Context, projectID string, params dto.ProjectRankingQueryParams, orgID string) (*dto.ProjectTaskRankingResponse, error) {
	return &dto.ProjectTaskRankingResponse{
		ProjectID: projectID,
		Rankings: []dto.TaskRankingItem{
			{
				Rank:          1,
				TaskID:        "task-1",
				Title:         "Important Task",
				PriorityScore: 95,
				Status:        "ready",
			},
			{
				Rank:          2,
				TaskID:        "task-2",
				Title:         "Another Task",
				PriorityScore: 80,
				Status:        "in_progress",
				AssigneeID:    strPtr("user-1"),
			},
		},
		Total: 2,
	}, nil
}

// RecalculatePriorities triggers dummy recalculation
func (s *DummyPriorityService) RecalculatePriorities(ctx context.Context, req dto.RecalculatePriorityRequest, orgID string) (*dto.RecalculatePrioritySyncResponse, *dto.RecalculatePriorityAsyncResponse, error) {
	if req.Async {
		return nil, &dto.RecalculatePriorityAsyncResponse{
			JobID:             "job_priority_abc123",
			Status:            "queued",
			EstimatedDuration: "5s",
		}, nil
	}
	return &dto.RecalculatePrioritySyncResponse{
		Recalculated:  12,
		Duration:      "150ms",
		AffectedTasks: []string{"task-1", "task-2"},
	}, nil, nil
}

// GetTaskRank returns a dummy rank
func (s *DummyPriorityService) GetTaskRank(ctx context.Context, taskID string, orgID string) (int, error) {
	return 1, nil
}

func strPtr(s string) *string {
	return &s
}
