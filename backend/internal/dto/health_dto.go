package dto

import "time"

// ===== Health Module DTOs =====

// ProjectHealthSummary represents a project summary in portfolio health
type ProjectHealthSummary struct {
	ProjectID  string `json:"projectId"`
	Name       string `json:"name"`
	HealthScore int  `json:"healthScore"`
	Status     string `json:"status"`
	Priority   int    `json:"priority"`
	Progress   int    `json:"progress"`
	Trend      string `json:"trend"`
}

// PortfolioHealthSummary represents the summary section of portfolio health
type PortfolioHealthSummary struct {
	TotalProjects int `json:"totalProjects"`
	Healthy       int `json:"healthy"`
	Caution       int `json:"caution"`
	AtRisk        int `json:"atRisk"`
	Critical      int `json:"critical"`
}

// PortfolioHealthResponse represents the portfolio health overview
type PortfolioHealthResponse struct {
	PortfolioHealthScore int                    `json:"portfolioHealthScore"`
	Status               string                 `json:"status"`
	Summary              PortfolioHealthSummary `json:"summary"`
	Projects             []ProjectHealthSummary `json:"projects"`
	CalculatedAt         time.Time              `json:"calculatedAt"`
}

// HealthBreakdown represents the breakdown of health score components
type HealthBreakdown struct {
	ScheduleHealth     int `json:"scheduleHealth"`
	CompletionHealth   int `json:"completionHealth"`
	DependencyHealth   int `json:"dependencyHealth"`
	ResourceHealth     int `json:"resourceHealth"`
	CriticalPathHealth int `json:"criticalPathHealth"`
}

// ScheduleDetails represents schedule health details
type ScheduleDetails struct {
	ExpectedProgress  int `json:"expectedProgress"`
	ActualProgress    int `json:"actualProgress"`
	Variance          int `json:"variance"`
	DaysUntilDeadline int `json:"daysUntilDeadline"`
}

// CompletionDetails represents completion health details
type CompletionDetails struct {
	TotalTasks     int `json:"totalTasks"`
	Completed      int `json:"completed"`
	InProgress     int `json:"inProgress"`
	CompletionRate int `json:"completionRate"`
}

// DependencyDetails represents dependency health details
type DependencyDetails struct {
	Total   int `json:"total"`
	Blocked int `json:"blocked"`
	AtRisk  int `json:"atRisk"`
}

// ResourceDetails represents resource health details
type ResourceDetails struct {
	TeamSize        int `json:"teamSize"`
	AvgAllocation   int `json:"avgAllocation"`
	Overallocated   int `json:"overallocated"`
	Underutilized   int `json:"underutilized"`
}

// HealthDetails represents detailed health information
type HealthDetails struct {
	Schedule    ScheduleDetails    `json:"schedule"`
	Completion  CompletionDetails  `json:"completion"`
	Dependencies DependencyDetails `json:"dependencies"`
	Resources   ResourceDetails    `json:"resources"`
}

// HealthTrend represents the trend information
type HealthTrend struct {
	Direction    string `json:"direction"`
	Change       int    `json:"change"`
	LastWeekScore int   `json:"lastWeekScore"`
}

// ProjectHealthResponse represents detailed project health
type ProjectHealthResponse struct {
	ProjectID    string          `json:"projectId"`
	ProjectName  string          `json:"projectName"`
	HealthScore  int             `json:"healthScore"`
	Status       string          `json:"status"`
	Breakdown    HealthBreakdown `json:"breakdown"`
	Details      HealthDetails   `json:"details"`
	Trend        HealthTrend     `json:"trend"`
	CalculatedAt time.Time       `json:"calculatedAt"`
}

// HealthDatapoint represents a single health data point over time
type HealthDatapoint struct {
	Date             string `json:"date"`
	HealthScore      int    `json:"healthScore"`
	ScheduleHealth   int    `json:"scheduleHealth"`
	CompletionHealth int    `json:"completionHealth"`
}

// HealthTrendPrediction represents health trend prediction
type HealthTrendPrediction struct {
	DaysUntilCritical int     `json:"daysUntilCritical"`
	Confidence        float64 `json:"confidence"`
}

// HealthTrendAnalysis represents the trend analysis
type HealthTrendAnalysis struct {
	Slope       float64               `json:"slope"`
	Direction   string                `json:"direction"`
	Prediction  HealthTrendPrediction `json:"prediction"`
}

// HealthTrendsResponse represents health trends over time
type HealthTrendsResponse struct {
	ProjectID  string              `json:"projectId"`
	TimeRange  string              `json:"timeRange"`
	Datapoints []HealthDatapoint   `json:"datapoints"`
	Trend      HealthTrendAnalysis `json:"trend"`
}

// HealthTrendsQueryParams represents query parameters for health trends
type HealthTrendsQueryParams struct {
	ProjectID string `form:"projectId" binding:"required"`
	Days      int    `form:"days,default=30" binding:"min=1,max=365"`
}

// ProjectHealthQueryParams represents query parameters for project health
type ProjectHealthQueryParams struct {
	IncludeBreakdown bool `form:"includeBreakdown"`
}
