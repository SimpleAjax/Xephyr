package dto

import "time"

// ===== Dependency Module DTOs =====

// DependencyInfo represents a direct dependency information
type DependencyInfo struct {
	DependencyID   string `json:"dependencyId"`
	DependsOnTaskID string `json:"dependsOnTaskId"`
	DependencyType string `json:"dependencyType"`
	LagHours       int    `json:"lagHours"`
	Status         string `json:"status"`
	IsBlocking     bool   `json:"isBlocking"`
}

// IndirectDependency represents an indirect dependency path
type IndirectDependency struct {
	Path  []string `json:"path"`
	Depth int      `json:"depth"`
}

// DependentInfo represents a task that depends on the current task
type DependentInfo struct {
	TaskID         string `json:"taskId"`
	DependencyType string `json:"dependencyType"`
	IsBlocked      bool   `json:"isBlocked"`
}

// ChainAnalysis represents dependency chain analysis
type ChainAnalysis struct {
	LongestChain       int    `json:"longestChain"`
	CriticalPathPosition string `json:"criticalPathPosition"`
	FloatHours         int    `json:"floatHours"`
}

// TaskDependenciesResponse represents task dependencies response
type TaskDependenciesResponse struct {
	TaskID        string                    `json:"taskId"`
	Dependencies  DependencySection         `json:"dependencies"`
	Dependents    DependentsSection         `json:"dependents"`
	ChainAnalysis ChainAnalysis             `json:"chainAnalysis"`
}

// DependencySection represents dependencies (direct and indirect)
type DependencySection struct {
	Direct   []DependencyInfo     `json:"direct"`
	Indirect []IndirectDependency `json:"indirect"`
}

// DependentsSection represents dependents (direct and indirect)
type DependentsSection struct {
	Direct   []DependentInfo      `json:"direct"`
	Indirect []IndirectDependency `json:"indirect"`
}

// CreateDependencyRequest represents a request to create a dependency
type CreateDependencyRequest struct {
	TaskID          string `json:"taskId" binding:"required"`
	DependsOnTaskID string `json:"dependsOnTaskId" binding:"required"`
	DependencyType  string `json:"dependencyType" binding:"required,oneof=finish_to_start start_to_start finish_to_finish start_to_finish"`
	LagHours        int    `json:"lagHours" binding:"min=0"`
}

// DependencyValidation represents validation result for a dependency
type DependencyValidation struct {
	Valid            bool   `json:"valid"`
	WouldCreateCycle bool   `json:"wouldCreateCycle"`
}

// DependencyImpact represents the impact of creating a dependency
type DependencyImpact struct {
	CriticalPathChanged bool      `json:"criticalPathChanged"`
	AffectedTasks       []string  `json:"affectedTasks"`
	NewProjectEndDate   *time.Time `json:"newProjectEndDate,omitempty"`
}

// CreateDependencyResponse represents the response after creating a dependency
type CreateDependencyResponse struct {
	DependencyID   string             `json:"dependencyId"`
	TaskID         string             `json:"taskId"`
	DependsOnTaskID string            `json:"dependsOnTaskId"`
	DependencyType string             `json:"dependencyType"`
	LagHours       int                `json:"lagHours"`
	CreatedAt      time.Time          `json:"createdAt"`
	Validation     DependencyValidation `json:"validation"`
	Impact         DependencyImpact   `json:"impact"`
}

// CriticalPathTask represents a task on the critical path
type CriticalPathTask struct {
	TaskID         string     `json:"taskId"`
	Title          string     `json:"title"`
	EstimatedHours int        `json:"estimatedHours"`
	EarliestStart  *time.Time `json:"earliestStart,omitempty"`
	EarliestFinish *time.Time `json:"earliestFinish,omitempty"`
	LatestStart    *time.Time `json:"latestStart,omitempty"`
	LatestFinish   *time.Time `json:"latestFinish,omitempty"`
	FloatHours     int        `json:"floatHours"`
}

// CriticalPathInfo represents the critical path information
type CriticalPathInfo struct {
	TaskIDs       []string           `json:"taskIds"`
	Tasks         []CriticalPathTask `json:"tasks"`
	TotalDuration int                `json:"totalDuration"`
}

// NonCriticalTask represents a task not on the critical path
type NonCriticalTask struct {
	TaskID     string `json:"taskId"`
	Title      string `json:"title"`
	FloatHours int    `json:"floatHours"`
}

// CriticalPathResponse represents the critical path response
type CriticalPathResponse struct {
	ProjectID         string            `json:"projectId"`
	CriticalPath      CriticalPathInfo  `json:"criticalPath"`
	NonCriticalTasks  []NonCriticalTask `json:"nonCriticalTasks"`
	ProjectDuration   int               `json:"projectDuration"`
	CalculatedAt      time.Time         `json:"calculatedAt"`
}

// ValidateDependencyRequest represents a request to validate a dependency
type ValidateDependencyRequest struct {
	TaskID          string `json:"taskId" binding:"required"`
	DependsOnTaskID string `json:"dependsOnTaskId" binding:"required"`
	DependencyType  string `json:"dependencyType" binding:"required,oneof=finish_to_start start_to_start finish_to_finish start_to_finish"`
}

// DependencyWarning represents a warning during validation
type DependencyWarning struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// ValidationImpact represents the estimated impact of a dependency
type ValidationImpact struct {
	EstimatedDelay int      `json:"estimatedDelay"`
	AffectedTasks  int      `json:"affectedTasks"`
}

// ValidateDependencyResponse represents the validation response
type ValidateDependencyResponse struct {
	Valid            bool                `json:"valid"`
	WouldCreateCycle bool                `json:"wouldCreateCycle"`
	Warnings         []DependencyWarning `json:"warnings"`
	Impact           ValidationImpact    `json:"impact"`
}

// DependencyGraphNode represents a node in the dependency graph
type DependencyGraphNode struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

// DependencyGraphEdge represents an edge in the dependency graph
type DependencyGraphEdge struct {
	ID       string `json:"id"`
	Source   string `json:"source"`
	Target   string `json:"target"`
	Type     string `json:"type"`
	LagHours int    `json:"lagHours"`
}

// DependencyGraphResponse represents the dependency graph response
type DependencyGraphResponse struct {
	ProjectID string                `json:"projectId"`
	Nodes     []DependencyGraphNode `json:"nodes"`
	Edges     []DependencyGraphEdge `json:"edges"`
}

// TaskDependenciesQueryParams represents query parameters for task dependencies
type TaskDependenciesQueryParams struct {
	IncludeIndirect bool `form:"includeIndirect"`
}

// CircularDependencyErrorDetails represents error details for circular dependency
type CircularDependencyErrorDetails struct {
	Cycle                 []string            `json:"cycle"`
	ExistingDependencies  []ExistingDependency `json:"existingDependencies"`
}

// ExistingDependency represents an existing dependency
type ExistingDependency struct {
	From string `json:"from"`
	To   string `json:"to"`
}
