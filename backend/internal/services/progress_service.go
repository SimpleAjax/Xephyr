package services

import (
	"context"
	"time"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
)

// ProgressService defines the interface for progress-related operations
type ProgressService interface {
	// GetProjectProgress returns progress for a project
	GetProjectProgress(ctx context.Context, projectID string, orgID string) (*dto.ProjectProgressResponse, error)

	// GetTaskProgress returns detailed progress for a task
	GetTaskProgress(ctx context.Context, taskID string, orgID string) (*dto.TaskProgressResponse, error)

	// UpdateTaskProgress updates the progress of a task
	UpdateTaskProgress(ctx context.Context, taskID string, req dto.UpdateTaskProgressRequest, orgID string, userID string) (*dto.TaskProgressUpdateResponse, error)

	// GetProjectRollup returns hierarchical progress rollup
	GetProjectRollup(ctx context.Context, projectID string, orgID string) (*dto.ProjectRollupResponse, error)

	// RecalculateProjectProgress recalculates project progress
	RecalculateProjectProgress(ctx context.Context, projectID string, orgID string) error
}

// DummyProgressService is a placeholder implementation of ProgressService
type DummyProgressService struct{}

// NewDummyProgressService creates a new dummy progress service
func NewDummyProgressService() ProgressService {
	return &DummyProgressService{}
}

// GetProjectProgress returns dummy project progress
func (s *DummyProgressService) GetProjectProgress(ctx context.Context, projectID string, orgID string) (*dto.ProjectProgressResponse, error) {
	return &dto.ProjectProgressResponse{
		ProjectID:          projectID,
		ProgressPercentage: 45,
		CalculationMethod:  "weighted_hours",
		Breakdown: dto.ProgressBreakdown{
			ByStatus: dto.ByStatusBreakdown{
				Backlog:    dto.StatusBreakdown{Count: 2, Hours: 80, Percentage: 20},
				Ready:      dto.StatusBreakdown{Count: 1, Hours: 50, Percentage: 12},
				InProgress: dto.StatusBreakdown{Count: 3, Hours: 120, Percentage: 30},
				Review:     dto.StatusBreakdown{Count: 1, Hours: 35, Percentage: 9},
				Done:       dto.StatusBreakdown{Count: 5, Hours: 115, Percentage: 29},
			},
			ByHierarchy: dto.ByHierarchyBreakdown{
				Tasks:    dto.HierarchyItem{Total: 8, Completed: 4},
				Subtasks: dto.HierarchyItem{Total: 4, Completed: 1},
			},
		},
		Variance: dto.ProgressVariance{
			ExpectedProgress: 50,
			ActualProgress:   45,
			Variance:         -5,
			Status:           "behind_schedule",
		},
		Milestones: []dto.Milestone{
			{
				TaskID:      "task-ec-1",
				Title:       "Design System Architecture",
				Status:      "completed",
				CompletedAt: timePtr(time.Now().UTC().Add(-30 * 24 * time.Hour)),
			},
		},
		CalculatedAt: time.Now().UTC(),
	}, nil
}

// GetTaskProgress returns dummy task progress
func (s *DummyProgressService) GetTaskProgress(ctx context.Context, taskID string, orgID string) (*dto.TaskProgressResponse, error) {
	return &dto.TaskProgressResponse{
		TaskID:             taskID,
		Title:              "Sample Task",
		Status:             "in_progress",
		ProgressPercentage: 60,
		EstimatedHours:     80,
		ActualHours:        48,
		RemainingHours:     32,
	}, nil
}

// UpdateTaskProgress updates dummy progress
func (s *DummyProgressService) UpdateTaskProgress(ctx context.Context, taskID string, req dto.UpdateTaskProgressRequest, orgID string, userID string) (*dto.TaskProgressUpdateResponse, error) {
	progress := 0
	if req.ProgressPercentage != nil {
		progress = *req.ProgressPercentage
	}
	return &dto.TaskProgressUpdateResponse{
		TaskID:             taskID,
		PreviousStatus:     "ready",
		NewStatus:          req.Status,
		ProgressPercentage: progress,
		ActualHours:        req.ActualHours,
		EstimatedHours:     80,
		RemainingHours:     32,
		Affected: dto.ProgressUpdateImpact{
			ParentProgressUpdated:  true,
			ProjectProgressUpdated: true,
			DependentsNotified:     []string{"task-ec-4"},
		},
	}, nil
}

// GetProjectRollup returns dummy rollup
func (s *DummyProgressService) GetProjectRollup(ctx context.Context, projectID string, orgID string) (*dto.ProjectRollupResponse, error) {
	return &dto.ProjectRollupResponse{
		ProjectID:          projectID,
		ProgressPercentage: 45,
		Tasks: []dto.RollupItem{
			{
				TaskID:             "task-1",
				Title:              "Parent Task",
				ProgressPercentage: 50,
				Status:             "in_progress",
				EstimatedHours:     100,
				ActualHours:        50,
				Children: []dto.RollupItem{
					{
						TaskID:             "task-1-1",
						Title:              "Subtask 1",
						ProgressPercentage: 100,
						Status:             "done",
						EstimatedHours:     40,
						ActualHours:        40,
					},
				},
			},
		},
		CalculatedAt: time.Now().UTC(),
	}, nil
}

// RecalculateProjectProgress is a no-op for dummy service
func (s *DummyProgressService) RecalculateProjectProgress(ctx context.Context, projectID string, orgID string) error {
	return nil
}
