package dto

// ===== Workload Module DTOs =====

// TaskAllocation represents a task in workload allocation
type TaskAllocation struct {
	TaskID             string  `json:"taskId"`
	Title              string  `json:"title"`
	ProjectID          string  `json:"projectId"`
	EstimatedHours     float64 `json:"estimatedHours"`
	AllocationThisWeek float64 `json:"allocationThisWeek"`
}

// MemberAllocation represents allocation details for a team member
type MemberAllocation struct {
	Percentage      int     `json:"percentage"`
	AssignedHours   float64 `json:"assignedHours"`
	CapacityHours   float64 `json:"capacityHours"`
}

// AvailabilityWindow represents availability for a time period
type AvailabilityWindow struct {
	ThisWeek float64 `json:"thisWeek"`
	NextWeek float64 `json:"nextWeek"`
}

// TeamMemberWorkload represents workload for a single team member
type TeamMemberWorkload struct {
	PersonID     string             `json:"personId"`
	Name         string             `json:"name"`
	Role         string             `json:"role"`
	Allocation   MemberAllocation   `json:"allocation"`
	Tasks        []TaskAllocation   `json:"tasks"`
	Status       string             `json:"status"`
	RiskLevel    string             `json:"riskLevel"`
	Availability AvailabilityWindow `json:"availability"`
}

// WorkloadSummary represents team workload summary
type WorkloadSummary struct {
	Overallocated  int `json:"overallocated"`
	Optimal        int `json:"optimal"`
	Available      int `json:"available"`
	Underutilized  int `json:"underutilized"`
}

// TeamWorkloadResponse represents team workload response
type TeamWorkloadResponse struct {
	WeekStarting    string             `json:"weekStarting"`
	TeamCapacity    float64            `json:"teamCapacity"`
	TeamAllocation  float64            `json:"teamAllocation"`
	UtilizationRate float64            `json:"utilizationRate"`
	Members         []TeamMemberWorkload `json:"members"`
	Summary         WorkloadSummary    `json:"summary"`
}

// TeamWorkloadQueryParams represents query parameters for team workload
type TeamWorkloadQueryParams struct {
	Week            string `form:"week"`
	IncludeForecast bool   `form:"includeForecast"`
}

// WeeklyForecast represents workload forecast for a week
type WeeklyForecast struct {
	WeekStarting string  `json:"weekStarting"`
	Allocation   int     `json:"allocation"`
	AssignedHours float64 `json:"assignedHours"`
	Tasks        int     `json:"tasks"`
	Risk         string  `json:"risk"`
}

// RiskPeriod represents a period of risk
type RiskPeriod struct {
	StartWeek string `json:"startWeek"`
	EndWeek   string `json:"endWeek"`
	Severity  string `json:"severity"`
	Reason    string `json:"reason"`
}

// WorkloadForecastResponse represents workload forecast response
type WorkloadForecastResponse struct {
	PersonID        string           `json:"personId"`
	Forecast        []WeeklyForecast `json:"forecast"`
	RiskPeriods     []RiskPeriod     `json:"riskPeriods"`
	Recommendations []string         `json:"recommendations"`
}

// WorkloadForecastQueryParams represents query parameters for forecast
type WorkloadForecastQueryParams struct {
	PersonID string `form:"personId" binding:"required"`
	Weeks    int    `form:"weeks,default=8" binding:"min=1,max=52"`
}

// UtilizationTrend represents utilization trend data point
type UtilizationTrend struct {
	Date         string  `json:"date"`
	Utilization  float64 `json:"utilization"`
	OptimalCount int     `json:"optimalCount"`
	OverCount    int     `json:"overCount"`
	UnderCount   int     `json:"underCount"`
}

// ProjectAllocation represents allocation by project
type ProjectAllocation struct {
	ProjectID    string  `json:"projectId"`
	ProjectName  string  `json:"projectName"`
	TotalHours   float64 `json:"totalHours"`
	Percentage   float64 `json:"percentage"`
}

// WorkloadDistribution represents workload distribution stats
type WorkloadDistribution struct {
	Overallocated   int     `json:"overallocated"`
	Optimal         int     `json:"optimal"`
	Available       int     `json:"available"`
	Underutilized   int     `json:"underutilized"`
}

// WorkloadAnalyticsResponse represents workload analytics response
type WorkloadAnalyticsResponse struct {
	Period              string               `json:"period"`
	AvgUtilization      float64              `json:"avgUtilization"`
	PeakUtilization     float64              `json:"peakUtilization"`
	Trends              []UtilizationTrend   `json:"trends"`
	ByProject           []ProjectAllocation  `json:"byProject"`
	Distribution        WorkloadDistribution `json:"distribution"`
}

// WorkloadAnalyticsQueryParams represents query parameters for analytics
type WorkloadAnalyticsQueryParams struct {
	Period string `form:"period,default=30d"`
}

// IndividualWorkloadResponse represents individual workload response
type IndividualWorkloadResponse struct {
	PersonID         string             `json:"personId"`
	Name             string             `json:"name"`
	CurrentAllocation MemberAllocation  `json:"currentAllocation"`
	Tasks            []TaskAllocation   `json:"tasks"`
	UpcomingWeeks    []WeeklyForecast   `json:"upcomingWeeks"`
	HistoricalTrend  []UtilizationTrend `json:"historicalTrend"`
}

// RebalanceWorkloadRequest represents a request to rebalance workload
type RebalanceWorkloadRequest struct {
	PersonID        string   `json:"personId,omitempty"`
	ProjectID       string   `json:"projectId,omitempty"`
	TargetDate      *string  `json:"targetDate,omitempty"`
	MaxUtilization  int      `json:"maxUtilization" binding:"min=80,max=120"`
}

// RebalanceSuggestion represents a rebalancing suggestion
type RebalanceSuggestion struct {
	TaskID       string  `json:"taskId"`
	TaskTitle    string  `json:"taskTitle"`
	CurrentOwner string  `json:"currentOwner"`
	SuggestedOwner string `json:"suggestedOwner"`
	Reason       string  `json:"reason"`
	Impact       string  `json:"impact"`
}

// RebalanceWorkloadResponse represents rebalancing suggestions response
type RebalanceWorkloadResponse struct {
	Suggestions []RebalanceSuggestion `json:"suggestions"`
	TotalImpact string                `json:"totalImpact"`
}
