package dto

import (
	"time"

	"github.com/google/uuid"
)

// ===== Scenario Module DTOs =====

// ProposedChanges represents proposed changes in a scenario
type ProposedChanges struct {
	PersonID         *string `json:"personId,omitempty"`
	LeaveStartDate   *string `json:"leaveStartDate,omitempty"`
	LeaveEndDate     *string `json:"leaveEndDate,omitempty"`
	CoverageStrategy *string `json:"coverageStrategy,omitempty"`
	// Additional change types can be added here
}

// CreateScenarioRequest represents a request to create a scenario
type CreateScenarioRequest struct {
	Title           string          `json:"title" binding:"required,max=200"`
	Description     string          `json:"description,omitempty"`
	ChangeType      string          `json:"changeType" binding:"required,oneof=employee_leave scope_change reallocation priority_shift"`
	ProposedChanges ProposedChanges `json:"proposedChanges" binding:"required"`
}

// ScenarioResponse represents a scenario response
type ScenarioResponse struct {
	ScenarioID       string                 `json:"scenarioId"`
	Title            string                 `json:"title"`
	ChangeType       string                 `json:"changeType"`
	Status           string                 `json:"status"`
	ProposedChanges  ProposedChanges        `json:"proposedChanges"`
	CreatedAt        time.Time              `json:"createdAt"`
	ImpactAnalysis   *ImpactAnalysisSummary `json:"impactAnalysis,omitempty"`
	SimulationStatus string                 `json:"simulationStatus"`
}

// ImpactAnalysisSummary represents a summary of impact analysis
type ImpactAnalysisSummary struct {
	TotalDelayDays   int     `json:"totalDelayDays"`
	CostImpact       float64 `json:"costImpact"`
	AffectedProjects int     `json:"affectedProjects"`
}

// ScenarioListQueryParams represents query parameters for listing scenarios
type ScenarioListQueryParams struct {
	Status string `form:"status"`
	Limit  int    `form:"limit,default=20" binding:"min=1,max=100"`
	Offset int    `form:"offset,default=0" binding:"min=0"`
}

// ScenarioListResponse represents the response for listing scenarios
type ScenarioListResponse struct {
	Scenarios []ScenarioResponse `json:"scenarios"`
	Total     int                `json:"total"`
}

// SimulateScenarioRequest represents a request to simulate a scenario
type SimulateScenarioRequest struct {
	Depth                  string `json:"depth" binding:"omitempty,oneof=quick full"`
	IncludeRecommendations bool   `json:"includeRecommendations"`
}

// AffectedProject represents a project affected by a scenario
type AffectedProject struct {
	ProjectID     string   `json:"projectId"`
	Name          string   `json:"name"`
	Impact        string   `json:"impact"`
	DelayDays     int      `json:"delayDays"`
	AffectedTasks []string `json:"affectedTasks"`
}

// AffectedTask represents a task affected by a scenario
type AffectedTask struct {
	TaskID                string                 `json:"taskId"`
	Title                 string                 `json:"title"`
	OriginalDueDate       *time.Time             `json:"originalDueDate"`
	NewDueDate            *time.Time             `json:"newDueDate"`
	DelayDays             int                    `json:"delayDays"`
	Reason                string                 `json:"reason"`
	SuggestedReassignment *SuggestedReassignment `json:"suggestedReassignment,omitempty"`
}

// SuggestedReassignment represents a suggested reassignment
type SuggestedReassignment struct {
	ToPersonID    string `json:"toPersonId"`
	Compatibility int    `json:"compatibility"`
}

// TimelineComparison represents timeline comparison
type TimelineComparison struct {
	OriginalEndDate *time.Time `json:"originalEndDate"`
	NewEndDate      *time.Time `json:"newEndDate"`
	TotalDelayDays  int        `json:"totalDelayDays"`
}

// CostBreakdownItem represents a cost breakdown item
type CostBreakdownItem struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

// CostAnalysis represents cost analysis
type CostAnalysis struct {
	TotalCost  float64             `json:"totalCost"`
	Breakdown  []CostBreakdownItem `json:"breakdown"`
	Confidence float64             `json:"confidence"`
}

// ResourceImpact represents resource impact
type ResourceImpact struct {
	PersonID          string `json:"personId"`
	CurrentAllocation int    `json:"currentAllocation"`
	NewAllocation     int    `json:"newAllocation"`
	Risk              string `json:"risk"`
}

// ImpactAnalysis represents detailed impact analysis
type ImpactAnalysis struct {
	AffectedProjects   []AffectedProject  `json:"affectedProjects"`
	AffectedTasks      []AffectedTask     `json:"affectedTasks"`
	TimelineComparison TimelineComparison `json:"timelineComparison"`
	CostAnalysis       CostAnalysis       `json:"costAnalysis"`
	ResourceImpacts    []ResourceImpact   `json:"resourceImpacts"`
}

// AIRecommendation represents an AI-generated recommendation
type AIRecommendation struct {
	Priority        int    `json:"priority"`
	Action          string `json:"action"`
	Reasoning       string `json:"reasoning"`
	EstimatedImpact string `json:"estimatedImpact"`
}

// SimulateScenarioResponse represents the simulation response
type SimulateScenarioResponse struct {
	ScenarioID         string             `json:"scenarioId"`
	SimulationStatus   string             `json:"simulationStatus"`
	ImpactAnalysis     ImpactAnalysis     `json:"impactAnalysis"`
	AIRecommendations  []AIRecommendation `json:"aiRecommendations"`
	CalculatedAt       time.Time          `json:"calculatedAt"`
	SimulationDuration string             `json:"simulationDuration"`
}

// ApplyScenarioRequest represents a request to apply a scenario
type ApplyScenarioRequest struct {
	ApplyRecommendations    bool `json:"applyRecommendations"`
	SelectedRecommendations []int `json:"selectedRecommendations,omitempty"`
	NotifyStakeholders      bool `json:"notifyStakeholders"`
}

// ReassignmentChange represents a reassignment change
type ReassignmentChange struct {
	TaskID string `json:"taskId"`
	From   string `json:"from"`
	To     string `json:"to"`
}

// DateAdjustment represents a date adjustment
type DateAdjustment struct {
	TaskID          string     `json:"taskId"`
	OriginalDueDate *time.Time `json:"originalDueDate"`
	NewDueDate      *time.Time `json:"newDueDate"`
}

// ScenarioChanges represents all changes applied
type ScenarioChanges struct {
	TasksReassigned   []ReassignmentChange `json:"tasksReassigned"`
	DatesAdjusted     []DateAdjustment     `json:"datesAdjusted"`
	NotificationsSent int                  `json:"notificationsSent"`
}

// ScenarioFollowUp represents follow-up actions
type ScenarioFollowUp struct {
	NudgesCreated         []string `json:"nudgesCreated"`
	CalendarEventsCreated bool     `json:"calendarEventsCreated"`
}

// ApplyScenarioResponse represents the response after applying a scenario
type ApplyScenarioResponse struct {
	ScenarioID string            `json:"scenarioId"`
	Status     string            `json:"status"`
	AppliedAt  time.Time         `json:"appliedAt"`
	AppliedBy  uuid.UUID         `json:"appliedBy"`
	Changes    ScenarioChanges   `json:"changes"`
	FollowUp   ScenarioFollowUp  `json:"followUp"`
}

// RejectScenarioResponse represents the response after rejecting a scenario
type RejectScenarioResponse struct {
	ScenarioID string    `json:"scenarioId"`
	Status     string    `json:"status"`
	RejectedAt time.Time `json:"rejectedAt"`
	RejectedBy uuid.UUID `json:"rejectedBy"`
}

// ModifyScenarioRequest represents a request to modify a scenario
type ModifyScenarioRequest struct {
	Title           *string          `json:"title,omitempty"`
	Description     *string          `json:"description,omitempty"`
	ProposedChanges *ProposedChanges `json:"proposedChanges,omitempty"`
}

// ScenarioHistoryEntry represents a history entry
type ScenarioHistoryEntry struct {
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	UserID    uuid.UUID `json:"userId"`
}

// ScenarioDetailResponse represents detailed scenario with history
type ScenarioDetailResponse struct {
	ScenarioResponse
	Description       string                 `json:"description"`
	ImpactAnalysis    *ImpactAnalysis        `json:"impactAnalysis,omitempty"`
	AIRecommendations []AIRecommendation     `json:"aiRecommendations,omitempty"`
	History           []ScenarioHistoryEntry `json:"history"`
}
