package dto

import "time"

// ===== Assignment Module DTOs =====

// PersonInfo represents basic person information
type PersonInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

// SkillMatchDetail represents detailed skill match information
type SkillMatchDetail struct {
	SkillID     string `json:"skillId"`
	Required    bool   `json:"required"`
	HasSkill    bool   `json:"hasSkill"`
	Proficiency int    `json:"proficiency"`
	MatchScore  int    `json:"matchScore"`
}

// CandidateBreakdown represents the breakdown of candidate score
type CandidateBreakdown struct {
	SkillMatch      int `json:"skillMatch"`
	Availability    int `json:"availability"`
	WorkloadBalance int `json:"workloadBalance"`
	PastPerformance int `json:"pastPerformance"`
}

// ContextSwitchAnalysis represents context switch analysis for a candidate
type ContextSwitchAnalysis struct {
	ActiveProjects int    `json:"activeProjects"`
	CurrentWorkload int   `json:"currentWorkload"`
	SwitchPenalty  int    `json:"switchPenalty"`
	RiskLevel      string `json:"riskLevel"`
}

// AssignmentCandidate represents a candidate for task assignment
type AssignmentCandidate struct {
	Rank                int                   `json:"rank"`
	Person              PersonInfo            `json:"person"`
	Score               int                   `json:"score"`
	Breakdown           CandidateBreakdown    `json:"breakdown"`
	SkillMatchDetails   []SkillMatchDetail    `json:"skillMatchDetails"`
	ContextSwitchAnalysis ContextSwitchAnalysis `json:"contextSwitchAnalysis"`
	Warnings            []string              `json:"warnings"`
	AIExplanation       string                `json:"aiExplanation"`
}

// AssignmentSuggestionsResponse represents assignment suggestions response
type AssignmentSuggestionsResponse struct {
	TaskID             string                `json:"taskId"`
	TaskTitle          string                `json:"taskTitle"`
	RequiredSkills     []string              `json:"requiredSkills"`
	Candidates         []AssignmentCandidate `json:"candidates"`
	UnassignableReason *string               `json:"unassignableReason,omitempty"`
}

// AssignmentSuggestionsQueryParams represents query parameters for suggestions
type AssignmentSuggestionsQueryParams struct {
	TaskID string `form:"taskId" binding:"required"`
	Limit  int    `form:"limit,default=3" binding:"min=1,max=10"`
}

// AssignTaskRequest represents a request to assign a task
type AssignTaskRequest struct {
	PersonID       string `json:"personId" binding:"required"`
	Note           string `json:"note,omitempty"`
	SkipSuggestion bool   `json:"skipSuggestion"`
}

// AssigneeInfo represents assignee information
type AssigneeInfo struct {
	PersonID string `json:"personId"`
	Name     string `json:"name"`
}

// AssignmentInfo represents assignment metadata
type AssignmentInfo struct {
	AssignedAt time.Time `json:"assignedAt"`
	AssignedBy string    `json:"assignedBy"`
}

// AssignmentImpact represents the impact of an assignment
type AssignmentImpact struct {
	WorkloadUpdated     bool     `json:"workloadUpdated"`
	NewAllocation       int      `json:"newAllocation"`
	NudgesGenerated     []string `json:"nudgesGenerated"`
	NotificationsSent   []string `json:"notificationsSent"`
}

// AssignTaskResponse represents the response after assigning a task
type AssignTaskResponse struct {
	TaskID           string           `json:"taskId"`
	AssignedTo       AssigneeInfo     `json:"assignedTo"`
	PreviousAssignee *AssigneeInfo    `json:"previousAssignee,omitempty"`
	Assignment       AssignmentInfo   `json:"assignment"`
	Impact           AssignmentImpact `json:"impact"`
}

// AutoAssignConstraints represents constraints for auto-assignment
type AutoAssignConstraints struct {
	MaxAllocation       int `json:"maxAllocation" binding:"min=0,max=200"`
	RequiredProficiency int `json:"requiredProficiency" binding:"min=1,max=4"`
}

// AutoAssignTaskRequest represents a request to auto-assign a task
type AutoAssignTaskRequest struct {
	Strategy    string                `json:"strategy" binding:"required,oneof=best_match balanced_workload fastest_completion"`
	Constraints AutoAssignConstraints `json:"constraints"`
}

// CompatibilityBreakdown represents compatibility score breakdown
type CompatibilityBreakdown struct {
	SkillMatch      int `json:"skillMatch"`
	Availability    int `json:"availability"`
	WorkloadBalance int `json:"workloadBalance"`
	PastPerformance int `json:"pastPerformance"`
}

// AssignmentCompatibilityResponse represents compatibility check response
type AssignmentCompatibilityResponse struct {
	TaskID       string                 `json:"taskId"`
	PersonID     string                 `json:"personId"`
	PersonName   string                 `json:"personName"`
	Score        int                    `json:"score"`
	Breakdown    CompatibilityBreakdown `json:"breakdown"`
	IsCompatible bool                   `json:"isCompatible"`
	Warnings     []string               `json:"warnings"`
	AIExplanation string                `json:"aiExplanation"`
}

// CompatibilityQueryParams represents query parameters for compatibility check
type CompatibilityQueryParams struct {
	TaskID   string `form:"taskId" binding:"required"`
	PersonID string `form:"personId" binding:"required"`
}

// ReassignmentItem represents a single reassignment in bulk operation
type ReassignmentItem struct {
	TaskID       string `json:"taskId" binding:"required"`
	FromPersonID string `json:"fromPersonId" binding:"required"`
	ToPersonID   string `json:"toPersonId" binding:"required"`
}

// BulkReassignRequest represents a request for bulk reassignment
type BulkReassignRequest struct {
	Reassignments []ReassignmentItem `json:"reassignments" binding:"required,min=1,max=50"`
	Reason        string             `json:"reason,omitempty"`
}

// ReassignmentResult represents the result of a single reassignment
type ReassignmentResult struct {
	TaskID   string `json:"taskId"`
	Status   string `json:"status"`
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Error    *ErrorInfo `json:"error,omitempty"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// BulkReassignResponse represents the response for bulk reassignment
type BulkReassignResponse struct {
	Processed   int                  `json:"processed"`
	Succeeded   int                  `json:"succeeded"`
	Failed      int                  `json:"failed"`
	Results     []ReassignmentResult `json:"results"`
	CompletedAt time.Time            `json:"completedAt"`
}
