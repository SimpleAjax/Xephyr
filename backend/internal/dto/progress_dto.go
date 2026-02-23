package dto

import "time"

// ===== Progress Module DTOs =====

// StatusBreakdown represents progress breakdown by status
type StatusBreakdown struct {
	Count      int     `json:"count"`
	Hours      float64 `json:"hours"`
	Percentage int     `json:"percentage"`
}

// ByStatusBreakdown represents progress breakdown by task status
type ByStatusBreakdown struct {
	Backlog    StatusBreakdown `json:"backlog"`
	Ready      StatusBreakdown `json:"ready"`
	InProgress StatusBreakdown `json:"in_progress"`
	Review     StatusBreakdown `json:"review"`
	Done       StatusBreakdown `json:"done"`
}

// ByHierarchyBreakdown represents progress breakdown by task hierarchy
type ByHierarchyBreakdown struct {
	Tasks    HierarchyItem `json:"tasks"`
	Subtasks HierarchyItem `json:"subtasks"`
}

// HierarchyItem represents tasks or subtasks count
type HierarchyItem struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
}

// ProgressVariance represents variance information
type ProgressVariance struct {
	ExpectedProgress int    `json:"expectedProgress"`
	ActualProgress   int    `json:"actualProgress"`
	Variance         int    `json:"variance"`
	Status           string `json:"status"`
}

// Milestone represents a project milestone
type Milestone struct {
	TaskID      string     `json:"taskId"`
	Title       string     `json:"title"`
	Status      string     `json:"status"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

// ProjectProgressResponse represents project progress response
type ProjectProgressResponse struct {
	ProjectID         string               `json:"projectId"`
	ProgressPercentage int                 `json:"progressPercentage"`
	CalculationMethod string              `json:"calculationMethod"`
	Breakdown         ProgressBreakdown   `json:"breakdown"`
	Variance          ProgressVariance    `json:"variance"`
	Milestones        []Milestone         `json:"milestones"`
	CalculatedAt      time.Time           `json:"calculatedAt"`
}

// ProgressBreakdown represents the complete progress breakdown
type ProgressBreakdown struct {
	ByStatus    ByStatusBreakdown    `json:"byStatus"`
	ByHierarchy ByHierarchyBreakdown `json:"byHierarchy"`
}

// TaskProgressResponse represents task progress details
type TaskProgressResponse struct {
	TaskID           string    `json:"taskId"`
	Title            string    `json:"title"`
	Status           string    `json:"status"`
	ProgressPercentage int     `json:"progressPercentage"`
	EstimatedHours   float64   `json:"estimatedHours"`
	ActualHours      float64   `json:"actualHours"`
	RemainingHours   float64   `json:"remainingHours"`
	StartDate        *time.Time `json:"startDate,omitempty"`
	DueDate          *time.Time `json:"dueDate,omitempty"`
	CompletedAt      *time.Time `json:"completedAt,omitempty"`
	SubtaskProgress  []SubtaskProgress `json:"subtaskProgress,omitempty"`
}

// SubtaskProgress represents subtask progress information
type SubtaskProgress struct {
	TaskID           string  `json:"taskId"`
	Title            string  `json:"title"`
	Status           string  `json:"status"`
	ProgressPercentage int   `json:"progressPercentage"`
}

// UpdateTaskProgressRequest represents a request to update task progress
type UpdateTaskProgressRequest struct {
	Status             string  `json:"status,omitempty" binding:"omitempty,oneof=backlog ready in_progress review done"`
	ProgressPercentage *int    `json:"progressPercentage,omitempty" binding:"omitempty,min=0,max=100"`
	ActualHours        float64 `json:"actualHours,omitempty" binding:"omitempty,min=0"`
	Note               string  `json:"note,omitempty"`
}

// TaskProgressUpdateResponse represents the response after updating progress
type TaskProgressUpdateResponse struct {
	TaskID           string              `json:"taskId"`
	PreviousStatus   string              `json:"previousStatus"`
	NewStatus        string              `json:"newStatus"`
	ProgressPercentage int               `json:"progressPercentage"`
	ActualHours      float64             `json:"actualHours"`
	EstimatedHours   float64             `json:"estimatedHours"`
	RemainingHours   float64             `json:"remainingHours"`
	Affected         ProgressUpdateImpact `json:"affected"`
}

// ProgressUpdateImpact represents the impact of a progress update
type ProgressUpdateImpact struct {
	ParentProgressUpdated bool     `json:"parentProgressUpdated"`
	ProjectProgressUpdated bool    `json:"projectProgressUpdated"`
	DependentsNotified    []string `json:"dependentsNotified"`
}

// RollupItem represents a hierarchical progress rollup item
type RollupItem struct {
	TaskID           string       `json:"taskId"`
	Title            string       `json:"title"`
	ProgressPercentage int        `json:"progressPercentage"`
	Status           string       `json:"status"`
	EstimatedHours   float64      `json:"estimatedHours"`
	ActualHours      float64      `json:"actualHours"`
	Children         []RollupItem `json:"children,omitempty"`
}

// ProjectRollupResponse represents hierarchical progress rollup for a project
type ProjectRollupResponse struct {
	ProjectID        string       `json:"projectId"`
	ProgressPercentage int        `json:"progressPercentage"`
	Tasks            []RollupItem `json:"tasks"`
	CalculatedAt     time.Time    `json:"calculatedAt"`
}
