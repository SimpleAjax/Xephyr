package dto

import "time"

// ===== Priority Module DTOs =====

// PriorityBreakdown represents the breakdown of priority score components
type PriorityBreakdown struct {
	ProjectPriority    int `json:"projectPriority"`
	BusinessValue      int `json:"businessValue"`
	DeadlineUrgency    int `json:"deadlineUrgency"`
	CriticalPathWeight int `json:"criticalPathWeight"`
	DependencyImpact   int `json:"dependencyImpact"`
}

// PriorityFactors represents the underlying factors for priority calculation
type PriorityFactors struct {
	ProjectPriority    int  `json:"projectPriority"`
	BusinessValue      int  `json:"businessValue"`
	DeadlineUrgency    int  `json:"deadlineUrgency"`
	IsOnCriticalPath   bool `json:"isOnCriticalPath"`
	BlockedTasksCount  int  `json:"blockedTasksCount"`
	DaysUntilDue       int  `json:"daysUntilDue"`
}

// TaskPriorityResponse represents a single task priority response
type TaskPriorityResponse struct {
	TaskID                 string            `json:"taskId"`
	UniversalPriorityScore int               `json:"universalPriorityScore"`
	Rank                   int               `json:"rank"`
	Breakdown              PriorityBreakdown `json:"breakdown"`
	Factors                PriorityFactors   `json:"factors"`
	CalculatedAt           time.Time         `json:"calculatedAt"`
}

// BulkPriorityRequest represents a request to get priorities for multiple tasks
type BulkPriorityRequest struct {
	TaskIds          []string `json:"taskIds" binding:"required,min=1,max=100"`
	IncludeBreakdown bool     `json:"includeBreakdown"`
}

// TaskPrioritySummary represents a simplified task priority for bulk responses
type TaskPrioritySummary struct {
	TaskID                 string `json:"taskId"`
	UniversalPriorityScore int    `json:"universalPriorityScore"`
	Rank                   int    `json:"rank"`
}

// BulkPriorityResponse represents the response for bulk priority request
type BulkPriorityResponse struct {
	Priorities  []TaskPrioritySummary `json:"priorities"`
	SortedOrder []string              `json:"sortedOrder"`
}

// TaskRankingItem represents a task in the project ranking
type TaskRankingItem struct {
	Rank          int     `json:"rank"`
	TaskID        string  `json:"taskId"`
	Title         string  `json:"title"`
	PriorityScore int     `json:"priorityScore"`
	Status        string  `json:"status"`
	AssigneeID    *string `json:"assigneeId,omitempty"`
}

// ProjectTaskRankingResponse represents the project task ranking response
type ProjectTaskRankingResponse struct {
	ProjectID string            `json:"projectId"`
	Rankings  []TaskRankingItem `json:"rankings"`
	Total     int               `json:"total"`
}

// RecalculatePriorityRequest represents a request to recalculate priorities
type RecalculatePriorityRequest struct {
	Scope     string  `json:"scope" binding:"required,oneof=project organization task"`
	ProjectID *string `json:"projectId,omitempty"`
	Async     bool    `json:"async"`
}

// RecalculatePrioritySyncResponse represents sync recalculation response
type RecalculatePrioritySyncResponse struct {
	Recalculated  int      `json:"recalculated"`
	Duration      string   `json:"duration"`
	AffectedTasks []string `json:"affectedTasks"`
}

// RecalculatePriorityAsyncResponse represents async recalculation response
type RecalculatePriorityAsyncResponse struct {
	JobID             string `json:"jobId"`
	Status            string `json:"status"`
	EstimatedDuration string `json:"estimatedDuration"`
}

// ProjectRankingQueryParams represents query parameters for project ranking
type ProjectRankingQueryParams struct {
	Status     string `form:"status"`
	AssigneeID string `form:"assigneeId"`
	MinScore   int    `form:"minScore,min=0,max=100"`
	Limit      int    `form:"limit,default=50" binding:"min=1,max=100"`
	Offset     int    `form:"offset,default=0" binding:"min=0"`
}
