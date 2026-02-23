package dto

import (
	"time"

	"github.com/google/uuid"
)

// ===== Nudge Module DTOs =====

// SuggestedAction represents an AI-suggested action for a nudge
type SuggestedAction struct {
	Type                string  `json:"type"`
	Description         string  `json:"description"`
	TargetTaskID        *string `json:"targetTaskId,omitempty"`
	SuggestedAssigneeID *string `json:"suggestedAssigneeId,omitempty"`
}

// RelatedEntities represents entities related to a nudge
type RelatedEntities struct {
	ProjectID *string `json:"projectId,omitempty"`
	TaskID    *string `json:"taskId,omitempty"`
	PersonID  *string `json:"personId,omitempty"`
}

// NudgeMetrics represents metrics associated with a nudge
type NudgeMetrics struct {
	AllocationPercentage int     `json:"allocationPercentage,omitempty"`
	AssignedTasks        int     `json:"assignedTasks,omitempty"`
	TotalHours           float64 `json:"totalHours,omitempty"`
}

// NudgeSummary represents nudge summary statistics
type NudgeSummary struct {
	Total       int            `json:"total"`
	Unread      int            `json:"unread"`
	BySeverity  map[string]int `json:"bySeverity"`
	ByType      map[string]int `json:"byType"`
}

// NudgeResponse represents a single nudge response
type NudgeResponse struct {
	ID              string          `json:"id"`
	Type            string          `json:"type"`
	Severity        string          `json:"severity"`
	Status          string          `json:"status"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	AIExplanation   string          `json:"aiExplanation"`
	SuggestedAction SuggestedAction `json:"suggestedAction"`
	RelatedEntities RelatedEntities `json:"relatedEntities"`
	Metrics         NudgeMetrics    `json:"metrics"`
	CriticalityScore int          `json:"criticalityScore"`
	CreatedAt       time.Time       `json:"createdAt"`
	ExpiresAt       *time.Time      `json:"expiresAt,omitempty"`
}

// NudgeListResponse represents the list of nudges response
type NudgeListResponse struct {
	Nudges  []NudgeResponse `json:"nudges"`
	Summary NudgeSummary    `json:"summary"`
}

// NudgeListQueryParams represents query parameters for listing nudges
type NudgeListQueryParams struct {
	Status    string `form:"status"`
	Severity  string `form:"severity"`
	Type      string `form:"type"`
	ProjectID string `form:"projectId"`
	PersonID  string `form:"personId"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Offset    int    `form:"offset,default=0" binding:"min=0"`
}

// NudgeActionRequest represents a request to take action on a nudge
type NudgeActionRequest struct {
	ActionType string                 `json:"actionType" binding:"required,oneof=accept_suggestion dismiss custom_action ask_alternatives snooze"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// NudgeActionResult represents the result of an action
type NudgeActionResult struct {
	TaskReassigned bool   `json:"taskReassigned"`
	FromUserID     string `json:"fromUserId,omitempty"`
	ToUserID       string `json:"toUserId,omitempty"`
	TaskID         string `json:"taskId,omitempty"`
}

// NudgeActionResponse represents the response after taking a nudge action
type NudgeActionResponse struct {
	NudgeID       string            `json:"nudgeId"`
	ActionTaken   string            `json:"actionTaken"`
	Result        NudgeActionResult `json:"result"`
	NudgeStatus   string            `json:"nudgeStatus"`
	FollowUpNudges []string         `json:"followUpNudges"`
	CompletedAt   time.Time         `json:"completedAt"`
}

// UpdateNudgeStatusRequest represents a request to update nudge status
type UpdateNudgeStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=unread read dismissed acted"`
}

// GenerateNudgesRequest represents a request to trigger nudge generation
type GenerateNudgesRequest struct {
	Scope     string   `json:"scope" binding:"required,oneof=project organization"`
	ProjectID *string  `json:"projectId,omitempty"`
	Types     []string `json:"types,omitempty"`
	Async     bool     `json:"async"`
}

// NudgeTypeStats represents statistics for a specific nudge type
type NudgeTypeStats struct {
	Generated int `json:"generated"`
	Acted     int `json:"acted"`
}

// NudgeStatsResponse represents nudge statistics response
type NudgeStatsResponse struct {
	Period          string                    `json:"period"`
	Generated       int                       `json:"generated"`
	Acted           int                       `json:"acted"`
	Dismissed       int                       `json:"dismissed"`
	Expired         int                       `json:"expired"`
	ActionRate      float64                   `json:"actionRate"`
	AvgTimeToAction string                    `json:"avgTimeToAction"`
	ByType          map[string]NudgeTypeStats `json:"byType"`
}

// NudgeStatsQueryParams represents query parameters for nudge stats
type NudgeStatsQueryParams struct {
	Period string `form:"period,default=30d"`
}

// NudgeAction represents a nudge action history entry
type NudgeAction struct {
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	UserID    uuid.UUID `json:"userId"`
}

// NudgeDetailResponse represents detailed nudge with history
type NudgeDetailResponse struct {
	NudgeResponse
	History []NudgeAction `json:"history"`
}
