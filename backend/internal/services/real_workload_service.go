package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/internal/repositories"
)

// RealWorkloadService implements WorkloadService using database queries
type RealWorkloadService struct {
	repos *repositories.Provider
}

// NewRealWorkloadService creates a new real workload service
func NewRealWorkloadService(repos *repositories.Provider) WorkloadService {
	return &RealWorkloadService{repos: repos}
}

// GetTeamWorkload returns team workload from database
func (s *RealWorkloadService) GetTeamWorkload(ctx context.Context, week string, includeForecast bool, orgID string) (*dto.TeamWorkloadResponse, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, err
	}

	// Parse week or use current week
	var weekStart time.Time
	if week != "" {
		weekStart, err = time.Parse("2006-01-02", week)
		if err != nil {
			weekStart = getCurrentWeekStart()
		}
	} else {
		weekStart = getCurrentWeekStart()
	}

	// Get team workload from repository
	tw, err := s.repos.GetWorkload().GetTeamWorkload(ctx, orgUUID, weekStart)
	if err != nil {
		return nil, err
	}

	// Get all users in organization for complete picture
	users, _, err := s.repos.GetUser().ListByOrganization(ctx, orgUUID, repositories.ListParams{Limit: 100})
	if err != nil {
		return nil, err
	}

	// Build member list
	members := make([]dto.TeamMemberWorkload, 0, len(users))
	overallocatedCount := 0
	optimalCount := 0
	availableCount := 0
	underutilizedCount := 0

	// Create a map of existing workload entries
	workloadMap := make(map[uuid.UUID]repositories.MemberWorkloadDetail)
	for _, mw := range tw.MemberWorkloads {
		workloadMap[mw.UserID] = mw
	}

	// Process all users
	for _, user := range users {
		// Check if user has workload entry
		if mw, exists := workloadMap[user.ID]; exists {
			status := s.getStatus(mw.AllocationPercentage)
			riskLevel := s.getRiskLevel(mw.AllocationPercentage)

			// Update counts
			switch status {
			case "overallocated":
				overallocatedCount++
			case "optimal":
				optimalCount++
			case "available":
				availableCount++
			case "underutilized":
				underutilizedCount++
			}

			// Get tasks for this user
			tasks, _, err := s.repos.GetTask().ListByAssignee(ctx, user.ID, repositories.ListParams{Limit: 50})
			if err != nil {
				tasks = []models.Task{}
			}

			taskAllocations := make([]dto.TaskAllocation, 0, len(tasks))
			for _, task := range tasks {
				taskAllocations = append(taskAllocations, dto.TaskAllocation{
					TaskID:             task.ID.String(),
					Title:              task.Title,
					ProjectID:          task.ProjectID.String(),
					EstimatedHours:     task.EstimatedHours,
					AllocationThisWeek: task.EstimatedHours, // Simplified
				})
			}

			members = append(members, dto.TeamMemberWorkload{
				PersonID: user.ID.String(),
				Name:     user.Name,
				Role:     s.getRole(user.ID, orgUUID),
				Allocation: dto.MemberAllocation{
					Percentage:    mw.AllocationPercentage,
					AssignedHours: mw.AssignedHours,
					CapacityHours: mw.CapacityHours,
				},
				Tasks:  taskAllocations,
				Status: status,
				Availability: dto.AvailabilityWindow{
					ThisWeek: mw.CapacityHours - mw.AssignedHours,
					NextWeek: mw.CapacityHours - mw.AssignedHours + 8, // Simplified forecast
				},
				RiskLevel: riskLevel,
			})
		} else {
			// User has no workload entry - they're available
			availableCount++
			members = append(members, dto.TeamMemberWorkload{
				PersonID: user.ID.String(),
				Name:     user.Name,
				Role:     s.getRole(user.ID, orgUUID),
				Allocation: dto.MemberAllocation{
					Percentage:    0,
					AssignedHours: 0,
					CapacityHours: 40,
				},
				Tasks:  []dto.TaskAllocation{},
				Status: "available",
				Availability: dto.AvailabilityWindow{
					ThisWeek: 40,
					NextWeek: 40,
				},
				RiskLevel: "low",
			})
		}
	}

	return &dto.TeamWorkloadResponse{
		WeekStarting:    weekStart.Format("2006-01-02"),
		TeamCapacity:    tw.TeamCapacity,
		TeamAllocation:  tw.TeamAllocation,
		UtilizationRate: tw.UtilizationRate,
		Members:         members,
		Summary: dto.WorkloadSummary{
			Overallocated: underutilizedCount,
			Optimal:       optimalCount,
			Available:     availableCount,
			Underutilized: underutilizedCount,
		},
	}, nil
}

// GetIndividualWorkload returns workload for a specific person
func (s *RealWorkloadService) GetIndividualWorkload(ctx context.Context, personID string, orgID string) (*dto.IndividualWorkloadResponse, error) {
	userUUID, err := uuid.Parse(personID)
	if err != nil {
		return nil, err
	}

	// Get user details
	user, err := s.repos.GetUser().GetByID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	// Get current week workload
	weekStart := getCurrentWeekStart()
	entry, err := s.repos.GetWorkload().GetByUserAndWeek(ctx, userUUID, weekStart)
	if err != nil {
		// No workload entry found, return empty
		return &dto.IndividualWorkloadResponse{
			PersonID: personID,
			Name:     user.Name,
			CurrentAllocation: dto.MemberAllocation{
				Percentage:    0,
				AssignedHours: 0,
				CapacityHours: 40,
			},
			Tasks:           []dto.TaskAllocation{},
			UpcomingWeeks:   []dto.WeeklyForecast{},
			HistoricalTrend: []dto.UtilizationTrend{},
		}, nil
	}

	// Get tasks assigned to user
	tasks, _, err := s.repos.GetTask().ListByAssignee(ctx, userUUID, repositories.ListParams{Limit: 50})
	if err != nil {
		tasks = []models.Task{}
	}

	taskAllocations := make([]dto.TaskAllocation, 0, len(tasks))
	for _, task := range tasks {
		taskAllocations = append(taskAllocations, dto.TaskAllocation{
			TaskID:             task.ID.String(),
			Title:              task.Title,
			ProjectID:          task.ProjectID.String(),
			EstimatedHours:     task.EstimatedHours,
			AllocationThisWeek: task.EstimatedHours,
		})
	}

	return &dto.IndividualWorkloadResponse{
		PersonID: personID,
		Name:     user.Name,
		CurrentAllocation: dto.MemberAllocation{
			Percentage:    entry.AllocationPercentage,
			AssignedHours: entry.TotalEstimatedHours,
			CapacityHours: entry.AvailableHours,
		},
		Tasks: taskAllocations,
		UpcomingWeeks: []dto.WeeklyForecast{
			{
				WeekStarting:  weekStart.Format("2006-01-02"),
				Allocation:    entry.AllocationPercentage,
				AssignedHours: entry.TotalEstimatedHours,
				Tasks:         entry.AssignedTasks,
				Risk:          s.getStatus(entry.AllocationPercentage),
			},
		},
		HistoricalTrend: []dto.UtilizationTrend{},
	}, nil
}

// GetWorkloadForecast returns workload forecast
func (s *RealWorkloadService) GetWorkloadForecast(ctx context.Context, personID string, weeks int, orgID string) (*dto.WorkloadForecastResponse, error) {
	userUUID, err := uuid.Parse(personID)
	if err != nil {
		return nil, err
	}

	_, err = s.repos.GetUser().GetByID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	// Get workload history
	now := time.Now()
	fromDate := now.AddDate(0, 0, -weeks*7)
	entries, err := s.repos.GetWorkload().ListByUser(ctx, userUUID, fromDate, now)
	if err != nil {
		entries = []models.WorkloadEntry{}
	}

	forecast := make([]dto.WeeklyForecast, 0, len(entries))
	for _, entry := range entries {
		forecast = append(forecast, dto.WeeklyForecast{
			WeekStarting:  entry.WeekStart.Format("2006-01-02"),
			Allocation:    entry.AllocationPercentage,
			AssignedHours: entry.TotalEstimatedHours,
			Tasks:         entry.AssignedTasks,
			Risk:          s.getStatus(entry.AllocationPercentage),
		})
	}

	return &dto.WorkloadForecastResponse{
		PersonID:        personID,
		Forecast:        forecast,
		RiskPeriods:     []dto.RiskPeriod{},
		Recommendations: []string{},
	}, nil
}

// GetWorkloadAnalytics returns workload analytics
func (s *RealWorkloadService) GetWorkloadAnalytics(ctx context.Context, period string, orgID string) (*dto.WorkloadAnalyticsResponse, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, err
	}

	weekStart := getCurrentWeekStart()
	tw, err := s.repos.GetWorkload().GetTeamWorkload(ctx, orgUUID, weekStart)
	if err != nil {
		return nil, err
	}

	// Calculate distribution
	overallocated := 0
	optimal := 0
	available := 0
	underutilized := 0

	for _, mw := range tw.MemberWorkloads {
		switch s.getStatus(mw.AllocationPercentage) {
		case "overallocated":
			overallocated++
		case "optimal":
			optimal++
		case "available":
			available++
		case "underutilized":
			underutilized++
		}
	}

	return &dto.WorkloadAnalyticsResponse{
		Period:          period,
		AvgUtilization:  tw.UtilizationRate,
		PeakUtilization: tw.UtilizationRate,
		Trends:          []dto.UtilizationTrend{},
		Distribution: dto.WorkloadDistribution{
			Overallocated: underutilized,
			Optimal:       optimal,
			Available:     available,
			Underutilized: underutilized,
		},
	}, nil
}

// GetRebalanceSuggestions returns workload rebalancing suggestions
func (s *RealWorkloadService) GetRebalanceSuggestions(ctx context.Context, req dto.RebalanceWorkloadRequest, orgID string) (*dto.RebalanceWorkloadResponse, error) {
	// This would implement actual rebalancing logic
	// For now, return empty suggestions
	return &dto.RebalanceWorkloadResponse{
		Suggestions: []dto.RebalanceSuggestion{},
		TotalImpact: "No rebalancing needed",
	}, nil
}

// Helper functions

func getCurrentWeekStart() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
}

func (s *RealWorkloadService) getStatus(percentage int) string {
	if percentage > 100 {
		return "overallocated"
	}
	if percentage >= 90 {
		return "optimal"
	}
	if percentage >= 50 {
		return "available"
	}
	return "underutilized"
}

func (s *RealWorkloadService) getRiskLevel(percentage int) string {
	if percentage > 120 {
		return "high"
	}
	if percentage > 100 {
		return "medium"
	}
	return "low"
}

func (s *RealWorkloadService) getRole(userID uuid.UUID, orgID uuid.UUID) string {
	// This would look up the user's role in the organization
	// For now, return a default
	return "Team Member"
}
