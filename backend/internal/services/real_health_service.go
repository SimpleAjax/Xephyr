package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/internal/repositories"
)

// RealHealthService implements HealthService using database queries
type RealHealthService struct {
	repos *repositories.Provider
}

// NewRealHealthService creates a new real health service
func NewRealHealthService(repos *repositories.Provider) HealthService {
	return &RealHealthService{repos: repos}
}

// GetPortfolioHealth returns portfolio health from database
func (s *RealHealthService) GetPortfolioHealth(ctx context.Context, orgID string) (*dto.PortfolioHealthResponse, error) {
	log.Printf("[RealHealthService] GetPortfolioHealth called with orgID: %s", orgID)
	
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		log.Printf("[RealHealthService] Error parsing orgID: %v", err)
		return nil, err
	}

	// Get all projects in organization
	projects, total, err := s.repos.GetProject().ListByOrganization(ctx, orgUUID, repositories.ListParams{Limit: 100})
	if err != nil {
		log.Printf("[RealHealthService] Error fetching projects: %v", err)
		return nil, err
	}
	
	log.Printf("[RealHealthService] Found %d projects (total: %d) for org %s", len(projects), total, orgID)

	// Calculate health metrics
	totalProjects := len(projects)
	healthyCount := 0
	cautionCount := 0
	atRiskCount := 0
	criticalCount := 0
	totalHealthScore := 0

	projectSummaries := make([]dto.ProjectHealthSummary, 0, len(projects))

	for _, project := range projects {
		healthScore := project.HealthScore
		totalHealthScore += healthScore

		status := s.getHealthStatus(healthScore)
		switch status {
		case "healthy":
			healthyCount++
		case "caution":
			cautionCount++
		case "at_risk":
			atRiskCount++
		case "critical":
			criticalCount++
		}

		projectSummaries = append(projectSummaries, dto.ProjectHealthSummary{
			ProjectID:   project.ID.String(),
			Name:        project.Name,
			HealthScore: healthScore,
			Status:      status,
			Priority:    project.Priority,
			Progress:    project.Progress,
			Trend:       s.calculateTrend(&project),
		})
	}

	portfolioScore := 0
	if totalProjects > 0 {
		portfolioScore = totalHealthScore / totalProjects
	}

	return &dto.PortfolioHealthResponse{
		PortfolioHealthScore: portfolioScore,
		Status:               s.getHealthStatus(portfolioScore),
		Summary: dto.PortfolioHealthSummary{
			TotalProjects: totalProjects,
			Healthy:       healthyCount,
			Caution:       cautionCount,
			AtRisk:        atRiskCount,
			Critical:      criticalCount,
		},
		Projects:     projectSummaries,
		CalculatedAt: time.Now().UTC(),
	}, nil
}

// GetProjectHealth returns detailed health for a specific project
func (s *RealHealthService) GetProjectHealth(ctx context.Context, projectID string, includeBreakdown bool, orgID string) (*dto.ProjectHealthResponse, error) {
	projUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, err
	}

	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, err
	}

	// Get project details
	project, err := s.repos.GetProject().GetByID(ctx, projUUID)
	if err != nil {
		return nil, err
	}

	// Get tasks for this project
	tasks, _, err := s.repos.GetTask().ListByProject(ctx, projUUID, repositories.ListParams{Limit: 1000})
	if err != nil {
		tasks = []models.Task{}
	}

	// Calculate completion metrics
	completedTasks := 0
	inProgressTasks := 0
	blockedTasks := 0
	atRiskTasks := 0
	for _, task := range tasks {
		switch task.Status {
		case models.TaskStatusDone:
			completedTasks++
		case models.TaskStatusInProgress:
			inProgressTasks++
		}
		if task.RiskScore > 50 {
			atRiskTasks++
		}
	}

	// Get team members count
	memberCount := len(project.Members)

	// Calculate overallocation
	overallocatedCount := 0
	weekStart := getCurrentWeekStart()
	workloadEntries, err := s.repos.GetWorkload().ListByOrganization(ctx, orgUUID, weekStart)
	if err == nil {
		for _, entry := range workloadEntries {
			if entry.AllocationPercentage > 100 {
				overallocatedCount++
			}
		}
	}

	totalTasks := len(tasks)
	completionRate := 0
	if totalTasks > 0 {
		completionRate = (completedTasks * 100) / totalTasks
	}

	resp := &dto.ProjectHealthResponse{
		ProjectID:   projectID,
		ProjectName: project.Name,
		HealthScore: project.HealthScore,
		Status:      s.getHealthStatus(project.HealthScore),
		Details: dto.HealthDetails{
			Schedule: dto.ScheduleDetails{
				ExpectedProgress:  s.calculateExpectedProgress(project),
				ActualProgress:    project.Progress,
				Variance:          project.Progress - s.calculateExpectedProgress(project),
				DaysUntilDeadline: s.daysUntilDeadline(project),
			},
			Completion: dto.CompletionDetails{
				TotalTasks:     totalTasks,
				Completed:      completedTasks,
				InProgress:     inProgressTasks,
				CompletionRate: completionRate,
			},
			Dependencies: dto.DependencyDetails{
				Total:   totalTasks, // Simplified
				Blocked: blockedTasks,
				AtRisk:  atRiskTasks,
			},
			Resources: dto.ResourceDetails{
				TeamSize:      memberCount,
				AvgAllocation: s.calculateAvgAllocation(workloadEntries),
				Overallocated: overallocatedCount,
				Underutilized: 0,
			},
		},
		Trend: dto.HealthTrend{
			Direction:     s.calculateTrend(project),
			Change:        0,
			LastWeekScore: project.HealthScore,
		},
		CalculatedAt: time.Now().UTC(),
	}

	if includeBreakdown {
		resp.Breakdown = dto.HealthBreakdown{
			ScheduleHealth:     s.calculateScheduleHealth(project),
			CompletionHealth:   s.calculateCompletionHealth(completionRate),
			DependencyHealth:   100 - (blockedTasks * 10), // Simplified
			ResourceHealth:     100 - (overallocatedCount * 20),
			CriticalPathHealth: project.HealthScore,
		}
	}

	return resp, nil
}

// GetBulkProjectHealth returns health for multiple projects
func (s *RealHealthService) GetBulkProjectHealth(ctx context.Context, projectIDs []string, orgID string) ([]dto.ProjectHealthSummary, error) {
	summaries := make([]dto.ProjectHealthSummary, 0, len(projectIDs))

	for _, pid := range projectIDs {
		health, err := s.GetProjectHealth(ctx, pid, false, orgID)
		if err != nil {
			continue
		}
		summaries = append(summaries, dto.ProjectHealthSummary{
			ProjectID:   health.ProjectID,
			Name:        health.ProjectName,
			HealthScore: health.HealthScore,
			Status:      health.Status,
		})
	}

	return summaries, nil
}

// GetHealthTrends returns health trends over time
func (s *RealHealthService) GetHealthTrends(ctx context.Context, projectID string, days int, orgID string) (*dto.HealthTrendsResponse, error) {
	projUUID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, err
	}

	project, err := s.repos.GetProject().GetByID(ctx, projUUID)
	if err != nil {
		return nil, err
	}

	// Generate mock trend data based on current health
	// In a real implementation, this would query a health_history table
	datapoints := []dto.HealthDatapoint{
		{
			Date:             time.Now().AddDate(0, 0, -days).Format("2006-01-02"),
			HealthScore:      project.HealthScore + 5,
			ScheduleHealth:   75,
			CompletionHealth: 60,
		},
		{
			Date:             time.Now().Format("2006-01-02"),
			HealthScore:      project.HealthScore,
			ScheduleHealth:   65,
			CompletionHealth: 45,
		},
	}

	return &dto.HealthTrendsResponse{
		ProjectID: projectID,
		TimeRange: "30d",
		Datapoints: datapoints,
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

// InvalidateHealthCache is a no-op for real service (caching not implemented)
func (s *RealHealthService) InvalidateHealthCache(ctx context.Context, projectID *string, orgID string) error {
	return nil
}

// Helper functions

func (s *RealHealthService) getHealthStatus(score int) string {
	if score >= 80 {
		return "healthy"
	}
	if score >= 60 {
		return "caution"
	}
	if score >= 40 {
		return "at_risk"
	}
	return "critical"
}

func (s *RealHealthService) calculateTrend(project *models.Project) string {
	if project == nil {
		return "stable"
	}
	// Simplified trend calculation based on health score
	if project.HealthScore >= 80 {
		return "stable"
	}
	if project.HealthScore >= 60 {
		return "stable"
	}
	return "worsening"
}

func (s *RealHealthService) calculateExpectedProgress(project *models.Project) int {
	if project == nil || project.StartDate == nil || project.TargetEndDate == nil {
		return 50
	}

	totalDuration := project.TargetEndDate.Sub(*project.StartDate).Hours()
	elapsed := time.Since(*project.StartDate).Hours()

	if totalDuration <= 0 {
		return 100
	}

	expected := int((elapsed / totalDuration) * 100)
	if expected > 100 {
		return 100
	}
	if expected < 0 {
		return 0
	}
	return expected
}

func (s *RealHealthService) daysUntilDeadline(project *models.Project) int {
	if project == nil || project.TargetEndDate == nil {
		return 30
	}
	days := int(project.TargetEndDate.Sub(time.Now()).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

func (s *RealHealthService) calculateAvgAllocation(entries []models.WorkloadEntry) int {
	if len(entries) == 0 {
		return 0
	}
	total := 0
	for _, e := range entries {
		total += e.AllocationPercentage
	}
	return total / len(entries)
}

func (s *RealHealthService) calculateScheduleHealth(project *models.Project) int {
	if project == nil {
		return 100
	}
	expected := s.calculateExpectedProgress(project)
	actual := project.Progress
	if expected == 0 {
		return 100
	}
	ratio := float64(actual) / float64(expected)
	if ratio >= 1.0 {
		return 100
	}
	return int(ratio * 100)
}

func (s *RealHealthService) calculateCompletionHealth(completionRate int) int {
	return completionRate
}
