package services

import (
	"context"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
)

// WorkloadService defines the interface for workload-related operations
type WorkloadService interface {
	// GetTeamWorkload returns team workload overview
	GetTeamWorkload(ctx context.Context, week string, includeForecast bool, orgID string) (*dto.TeamWorkloadResponse, error)

	// GetIndividualWorkload returns workload for a specific person
	GetIndividualWorkload(ctx context.Context, personID string, orgID string) (*dto.IndividualWorkloadResponse, error)

	// GetWorkloadForecast returns workload forecast
	GetWorkloadForecast(ctx context.Context, personID string, weeks int, orgID string) (*dto.WorkloadForecastResponse, error)

	// GetWorkloadAnalytics returns workload analytics
	GetWorkloadAnalytics(ctx context.Context, period string, orgID string) (*dto.WorkloadAnalyticsResponse, error)

	// GetRebalanceSuggestions returns workload rebalancing suggestions
	GetRebalanceSuggestions(ctx context.Context, req dto.RebalanceWorkloadRequest, orgID string) (*dto.RebalanceWorkloadResponse, error)
}

// DummyWorkloadService is a placeholder implementation of WorkloadService
type DummyWorkloadService struct{}

// NewDummyWorkloadService creates a new dummy workload service
func NewDummyWorkloadService() WorkloadService {
	return &DummyWorkloadService{}
}

// GetTeamWorkload returns dummy team workload
func (s *DummyWorkloadService) GetTeamWorkload(ctx context.Context, week string, includeForecast bool, orgID string) (*dto.TeamWorkloadResponse, error) {
	return &dto.TeamWorkloadResponse{
		WeekStarting:    "2026-02-24",
		TeamCapacity:    320,
		TeamAllocation:  385,
		UtilizationRate: 1.20,
		Members: []dto.TeamMemberWorkload{
			{
				PersonID: "user-emma",
				Name:     "Emma Wilson",
				Role:     "Designer",
				Allocation: dto.MemberAllocation{
					Percentage:    125,
					AssignedHours: 75,
					CapacityHours: 40,
				},
				Tasks: []dto.TaskAllocation{
					{
						TaskID:             "task-fit-1",
						Title:              "Mobile App UI Design",
						ProjectID:          "proj-mobile",
						EstimatedHours:     40,
						AllocationThisWeek: 40,
					},
					{
						TaskID:             "task-ec-5",
						Title:              "E-Commerce Admin Dashboard",
						ProjectID:          "proj-ecommerce",
						EstimatedHours:     25,
						AllocationThisWeek: 25,
					},
				},
				Status:    "overallocated",
				RiskLevel: "high",
				Availability: dto.AvailabilityWindow{
					ThisWeek: 0,
					NextWeek: 5,
				},
			},
			{
				PersonID: "user-mike",
				Name:     "Mike Rodriguez",
				Allocation: dto.MemberAllocation{
					Percentage:    40,
					AssignedHours: 16,
					CapacityHours: 40,
				},
				Status: "available",
				Availability: dto.AvailabilityWindow{
					ThisWeek: 24,
					NextWeek: 35,
				},
			},
		},
		Summary: dto.WorkloadSummary{
			Overallocated: 2,
			Optimal:       3,
			Available:     2,
			Underutilized: 1,
		},
	}, nil
}

// GetIndividualWorkload returns dummy individual workload
func (s *DummyWorkloadService) GetIndividualWorkload(ctx context.Context, personID string, orgID string) (*dto.IndividualWorkloadResponse, error) {
	return &dto.IndividualWorkloadResponse{
		PersonID: personID,
		Name:     "Emma Wilson",
		CurrentAllocation: dto.MemberAllocation{
			Percentage:    125,
			AssignedHours: 75,
			CapacityHours: 40,
		},
		Tasks: []dto.TaskAllocation{
			{
				TaskID:             "task-fit-1",
				Title:              "Mobile App UI Design",
				ProjectID:          "proj-mobile",
				EstimatedHours:     40,
				AllocationThisWeek: 40,
			},
		},
		UpcomingWeeks: []dto.WeeklyForecast{
			{
				WeekStarting:  "2026-02-24",
				Allocation:    125,
				AssignedHours: 75,
				Tasks:         4,
				Risk:          "overallocated",
			},
		},
		HistoricalTrend: []dto.UtilizationTrend{
			{
				Date:         "2026-02-01",
				Utilization:  1.1,
				OptimalCount: 5,
			},
		},
	}, nil
}

// GetWorkloadForecast returns dummy forecast
func (s *DummyWorkloadService) GetWorkloadForecast(ctx context.Context, personID string, weeks int, orgID string) (*dto.WorkloadForecastResponse, error) {
	return &dto.WorkloadForecastResponse{
		PersonID: personID,
		Forecast: []dto.WeeklyForecast{
			{
				WeekStarting:  "2026-02-24",
				Allocation:    125,
				AssignedHours: 75,
				Tasks:         4,
				Risk:          "overallocated",
			},
			{
				WeekStarting:  "2026-03-03",
				Allocation:    95,
				AssignedHours: 57,
				Tasks:         3,
				Risk:          "optimal",
			},
		},
		RiskPeriods: []dto.RiskPeriod{
			{
				StartWeek: "2026-02-24",
				EndWeek:   "2026-03-02",
				Severity:  "high",
				Reason:    "Overallocated at 125%",
			},
		},
		Recommendations: []string{
			"Consider reassigning 1 task from Feb 24 week",
		},
	}, nil
}

// GetWorkloadAnalytics returns dummy analytics
func (s *DummyWorkloadService) GetWorkloadAnalytics(ctx context.Context, period string, orgID string) (*dto.WorkloadAnalyticsResponse, error) {
	return &dto.WorkloadAnalyticsResponse{
		Period:         period,
		AvgUtilization: 0.85,
		PeakUtilization: 1.25,
		Trends: []dto.UtilizationTrend{
			{
				Date:         "2026-02-01",
				Utilization:  0.8,
				OptimalCount: 5,
				OverCount:    1,
				UnderCount:   0,
			},
		},
		ByProject: []dto.ProjectAllocation{
			{
				ProjectID:   "proj-ecommerce",
				ProjectName: "E-Commerce Platform",
				TotalHours:  200,
				Percentage:  0.4,
			},
		},
		Distribution: dto.WorkloadDistribution{
			Overallocated: 2,
			Optimal:       3,
			Available:     2,
			Underutilized: 1,
		},
	}, nil
}

// GetRebalanceSuggestions returns dummy suggestions
func (s *DummyWorkloadService) GetRebalanceSuggestions(ctx context.Context, req dto.RebalanceWorkloadRequest, orgID string) (*dto.RebalanceWorkloadResponse, error) {
	return &dto.RebalanceWorkloadResponse{
		Suggestions: []dto.RebalanceSuggestion{
			{
				TaskID:         "task-web-3",
				TaskTitle:      "Marketing Website",
				CurrentOwner:   "user-emma",
				SuggestedOwner: "user-rachel",
				Reason:         "Reduces Emma's allocation from 125% to 95%",
				Impact:         "Low risk, Rachel has required skills",
			},
		},
		TotalImpact: "Reduces team overallocation by 1 person",
	}, nil
}
