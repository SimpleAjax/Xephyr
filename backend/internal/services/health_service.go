package services

import (
	"context"
	"time"

	"github.com/SimpleAjax/Xephyr/internal/dto"
)

// HealthService defines the interface for health-related operations
type HealthService interface {
	// GetPortfolioHealth returns the overall portfolio health
	GetPortfolioHealth(ctx context.Context, orgID string) (*dto.PortfolioHealthResponse, error)

	// GetProjectHealth returns detailed health for a specific project
	GetProjectHealth(ctx context.Context, projectID string, includeBreakdown bool, orgID string) (*dto.ProjectHealthResponse, error)

	// GetBulkProjectHealth returns health for multiple projects
	GetBulkProjectHealth(ctx context.Context, projectIDs []string, orgID string) ([]dto.ProjectHealthSummary, error)

	// GetHealthTrends returns health trends over time
	GetHealthTrends(ctx context.Context, projectID string, days int, orgID string) (*dto.HealthTrendsResponse, error)

	// InvalidateHealthCache invalidates cached health data
	InvalidateHealthCache(ctx context.Context, projectID *string, orgID string) error
}

// DummyHealthService is a placeholder implementation of HealthService
type DummyHealthService struct{}

// NewDummyHealthService creates a new dummy health service
func NewDummyHealthService() HealthService {
	return &DummyHealthService{}
}

// GetPortfolioHealth returns dummy portfolio health
func (s *DummyHealthService) GetPortfolioHealth(ctx context.Context, orgID string) (*dto.PortfolioHealthResponse, error) {
	return &dto.PortfolioHealthResponse{
		PortfolioHealthScore: 68,
		Status:               "caution",
		Summary: dto.PortfolioHealthSummary{
			TotalProjects: 5,
			Healthy:       2,
			Caution:       2,
			AtRisk:        1,
			Critical:      0,
		},
		Projects: []dto.ProjectHealthSummary{
			{
				ProjectID:   "proj-ecommerce",
				Name:        "E-Commerce Platform",
				HealthScore: 72,
				Status:      "caution",
				Priority:    95,
				Progress:    45,
				Trend:       "stable",
			},
			{
				ProjectID:   "proj-mobile",
				Name:        "Fitness App",
				HealthScore: 45,
				Status:      "at_risk",
				Priority:    88,
				Progress:    25,
				Trend:       "worsening",
			},
		},
		CalculatedAt: time.Now().UTC(),
	}, nil
}

// GetProjectHealth returns dummy project health
func (s *DummyHealthService) GetProjectHealth(ctx context.Context, projectID string, includeBreakdown bool, orgID string) (*dto.ProjectHealthResponse, error) {
	resp := &dto.ProjectHealthResponse{
		ProjectID:   projectID,
		ProjectName: "Sample Project",
		HealthScore: 72,
		Status:      "caution",
		Details: dto.HealthDetails{
			Schedule: dto.ScheduleDetails{
				ExpectedProgress:  50,
				ActualProgress:    45,
				Variance:          -5,
				DaysUntilDeadline: 67,
			},
			Completion: dto.CompletionDetails{
				TotalTasks:     12,
				Completed:      5,
				InProgress:     4,
				CompletionRate: 45,
			},
			Dependencies: dto.DependencyDetails{
				Total:   8,
				Blocked: 1,
				AtRisk:  2,
			},
			Resources: dto.ResourceDetails{
				TeamSize:      4,
				AvgAllocation: 95,
				Overallocated: 1,
				Underutilized: 0,
			},
		},
		Trend: dto.HealthTrend{
			Direction:     "stable",
			Change:        0,
			LastWeekScore: 72,
		},
		CalculatedAt: time.Now().UTC(),
	}

	if includeBreakdown {
		resp.Breakdown = dto.HealthBreakdown{
			ScheduleHealth:     65,
			CompletionHealth:   45,
			DependencyHealth:   80,
			ResourceHealth:     85,
			CriticalPathHealth: 70,
		}
	}

	return resp, nil
}

// GetBulkProjectHealth returns dummy bulk health
func (s *DummyHealthService) GetBulkProjectHealth(ctx context.Context, projectIDs []string, orgID string) ([]dto.ProjectHealthSummary, error) {
	summaries := make([]dto.ProjectHealthSummary, len(projectIDs))
	for i, id := range projectIDs {
		summaries[i] = dto.ProjectHealthSummary{
			ProjectID:   id,
			Name:        "Project " + id,
			HealthScore: 70,
			Status:      "healthy",
			Priority:    80,
			Progress:    50,
			Trend:       "stable",
		}
	}
	return summaries, nil
}

// GetHealthTrends returns dummy health trends
func (s *DummyHealthService) GetHealthTrends(ctx context.Context, projectID string, days int, orgID string) (*dto.HealthTrendsResponse, error) {
	return &dto.HealthTrendsResponse{
		ProjectID: projectID,
		TimeRange: "30d",
		Datapoints: []dto.HealthDatapoint{
			{
				Date:             "2026-01-23",
				HealthScore:      78,
				ScheduleHealth:   75,
				CompletionHealth: 60,
			},
			{
				Date:             "2026-02-22",
				HealthScore:      72,
				ScheduleHealth:   65,
				CompletionHealth: 45,
			},
		},
		Trend: dto.HealthTrendAnalysis{
			Slope:     -0.2,
			Direction: "declining",
			Prediction: dto.HealthTrendPrediction{
				DaysUntilCritical: 45,
				Confidence:        0.75,
			},
		},
	}, nil
}

// InvalidateHealthCache is a no-op for dummy service
func (s *DummyHealthService) InvalidateHealthCache(ctx context.Context, projectID *string, orgID string) error {
	return nil
}
