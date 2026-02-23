package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

// BeforeCreate generates UUID before inserting
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return nil
}

// ===== User & Organization Models =====

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RolePM     UserRole = "pm"
	RoleMember UserRole = "member"
)

type User struct {
	BaseModel
	Email       string    `json:"email" gorm:"uniqueIndex;not null"`
	Name        string    `json:"name"`
	AvatarURL   string    `json:"avatarUrl"`
	PasswordHash string   `json:"-" gorm:"column:password_hash"`
	HourlyRate  float64   `json:"hourlyRate"`
	Timezone    string    `json:"timezone"`
	IsActive    bool      `json:"isActive" gorm:"default:true"`
	
	// Relationships
	OrganizationMemberships []OrganizationMember `json:"organizationMemberships,omitempty" gorm:"foreignKey:UserID"`
	Skills                  []UserSkill          `json:"skills,omitempty" gorm:"foreignKey:UserID"`
}

type Organization struct {
	BaseModel
	Name    string `json:"name"`
	Slug    string `json:"slug" gorm:"uniqueIndex"`
	Plan    string `json:"plan" gorm:"default:'free'"`
	Settings JSONB `json:"settings,omitempty" gorm:"type:jsonb"`
	
	// Relationships
	Members  []OrganizationMember `json:"members,omitempty" gorm:"foreignKey:OrganizationID"`
	Projects []Project            `json:"projects,omitempty" gorm:"foreignKey:OrganizationID"`
}

type OrganizationMember struct {
	BaseModel
	OrganizationID uuid.UUID    `json:"organizationId" gorm:"not null"`
	UserID         uuid.UUID    `json:"userId" gorm:"not null"`
	Role           UserRole     `json:"role" gorm:"default:'member'"`
	JoinedAt       time.Time    `json:"joinedAt" gorm:"default:CURRENT_TIMESTAMP"`
	
	Organization Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	User         User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ===== Skill Models =====

type Skill struct {
	BaseModel
	OrganizationID  *uuid.UUID `json:"organizationId,omitempty"` // nil for global skills
	Name            string     `json:"name"`
	Category        string     `json:"category"`
	Description     string     `json:"description"`
	
	UserSkills      []UserSkill `json:"userSkills,omitempty" gorm:"foreignKey:SkillID"`
	TaskSkills      []TaskSkill `json:"taskSkills,omitempty" gorm:"foreignKey:SkillID"`
}

type UserSkill struct {
	BaseModel
	UserID       uuid.UUID `json:"userId" gorm:"not null"`
	SkillID      uuid.UUID `json:"skillId" gorm:"not null"`
	Proficiency  int       `json:"proficiency" gorm:"check:proficiency >= 1 AND proficiency <= 4"` // 1-4 scale
	YearsOfExperience float64 `json:"yearsOfExperience"`
	
	User    User    `json:"-" gorm:"foreignKey:UserID"`
	Skill   Skill   `json:"skill,omitempty" gorm:"foreignKey:SkillID"`
}

type TaskSkill struct {
	BaseModel
	TaskID             uuid.UUID `json:"taskId" gorm:"not null"`
	SkillID            uuid.UUID `json:"skillId" gorm:"not null"`
	ProficiencyRequired int      `json:"proficiencyRequired" gorm:"default:3"` // min proficiency needed
	IsRequired         bool      `json:"isRequired" gorm:"default:true"`
	
	Task  Task  `json:"-" gorm:"foreignKey:TaskID"`
	Skill Skill `json:"skill,omitempty" gorm:"foreignKey:SkillID"`
}

// ===== Project Models =====

type ProjectStatus string

const (
	ProjectActive    ProjectStatus = "active"
	ProjectPaused    ProjectStatus = "paused"
	ProjectCompleted ProjectStatus = "completed"
	ProjectArchived  ProjectStatus = "archived"
)

type Project struct {
	BaseModel
	OrganizationID uuid.UUID     `json:"organizationId" gorm:"not null"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Status         ProjectStatus `json:"status" gorm:"default:'active'"`
	Priority       int           `json:"priority" gorm:"default:50"` // 0-100
	HealthScore    int           `json:"healthScore" gorm:"default:100"`
	Progress       int           `json:"progress" gorm:"default:0"` // 0-100
	StartDate      *time.Time    `json:"startDate"`
	TargetEndDate  *time.Time    `json:"targetEndDate"`
	Budget         float64       `json:"budget"`
	
	// Relationships
	Organization Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Tasks        []Task       `json:"tasks,omitempty" gorm:"foreignKey:ProjectID"`
	Members      []ProjectMember `json:"members,omitempty" gorm:"foreignKey:ProjectID"`
}

type ProjectMember struct {
	BaseModel
	ProjectID uuid.UUID `json:"projectId" gorm:"not null"`
	UserID    uuid.UUID `json:"userId" gorm:"not null"`
	Role      string    `json:"role"`
	
	Project Project `json:"-" gorm:"foreignKey:ProjectID"`
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ===== Task Models =====

type TaskStatus string

const (
	TaskStatusBacklog    TaskStatus = "backlog"
	TaskStatusReady      TaskStatus = "ready"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusReview     TaskStatus = "review"
	TaskStatusDone       TaskStatus = "done"
)

type TaskPriority string

const (
	TaskPriorityLow      TaskPriority = "low"
	TaskPriorityMedium   TaskPriority = "medium"
	TaskPriorityHigh     TaskPriority = "high"
	TaskPriorityCritical TaskPriority = "critical"
)

type Task struct {
	BaseModel
	ProjectID       uuid.UUID    `json:"projectId" gorm:"not null"`
	ParentTaskID    *uuid.UUID   `json:"parentTaskId,omitempty"`
	HierarchyLevel  int          `json:"hierarchyLevel" gorm:"default:1"` // 1=task, 2=subtask, 3=grandchild
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	Status          TaskStatus   `json:"status" gorm:"default:'backlog'"`
	Priority        TaskPriority `json:"priority" gorm:"default:'medium'"`
	PriorityScore   int          `json:"priorityScore" gorm:"default:0"` // 0-100 AI-calculated
	BusinessValue   int          `json:"businessValue" gorm:"default:50"` // 0-100
	EstimatedHours  float64      `json:"estimatedHours"`
	ActualHours     float64      `json:"actualHours"`
	StartDate       *time.Time   `json:"startDate"`
	DueDate         *time.Time   `json:"dueDate"`
	CompletedAt     *time.Time   `json:"completedAt"`
	AssigneeID      *uuid.UUID   `json:"assigneeId,omitempty"`
	IsMilestone     bool         `json:"isMilestone" gorm:"default:false"`
	IsCriticalPath  bool         `json:"isCriticalPath" gorm:"default:false"`
	RiskScore       int          `json:"riskScore" gorm:"default:0"`
	
	// Relationships
	Project      Project       `json:"-" gorm:"foreignKey:ProjectID"`
	Assignee     *User         `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
	ParentTask   *Task         `json:"-" gorm:"foreignKey:ParentTaskID"`
	Subtasks     []Task        `json:"subtasks,omitempty" gorm:"foreignKey:ParentTaskID"`
	Skills       []TaskSkill   `json:"skills,omitempty" gorm:"foreignKey:TaskID"`
	Dependencies []TaskDependency `json:"dependencies,omitempty" gorm:"foreignKey:TaskID"`
	BlockedBy    []TaskDependency `json:"blockedBy,omitempty" gorm:"foreignKey:DependsOnTaskID"`
}

// ===== Dependency Models =====

type DependencyType string

const (
	DependencyFinishToStart DependencyType = "finish_to_start"
	DependencyStartToStart  DependencyType = "start_to_start"
	DependencyFinishToFinish DependencyType = "finish_to_finish"
	DependencyStartToFinish DependencyType = "start_to_finish"
)

type TaskDependency struct {
	BaseModel
	TaskID           uuid.UUID      `json:"taskId" gorm:"not null"`
	DependsOnTaskID  uuid.UUID      `json:"dependsOnTaskId" gorm:"not null"`
	DependencyType   DependencyType `json:"dependencyType" gorm:"default:'finish_to_start'"`
	LagHours         int            `json:"lagHours" gorm:"default:0"`
	
	Task          Task `json:"-" gorm:"foreignKey:TaskID"`
	DependsOnTask Task `json:"-" gorm:"foreignKey:DependsOnTaskID"`
}

// ===== Nudge Models =====

type NudgeType string

const (
	NudgeTypeOverload         NudgeType = "overload"
	NudgeTypeDelayRisk        NudgeType = "delay_risk"
	NudgeTypeSkillGap         NudgeType = "skill_gap"
	NudgeTypeUnassigned       NudgeType = "unassigned"
	NudgeTypeBlocked          NudgeType = "blocked"
	NudgeTypeConflict         NudgeType = "conflict"
	NudgeTypeDependencyBlock  NudgeType = "dependency_block"
)

type NudgeSeverity string

const (
	NudgeSeverityLow    NudgeSeverity = "low"
	NudgeSeverityMedium NudgeSeverity = "medium"
	NudgeSeverityHigh   NudgeSeverity = "high"
)

type NudgeStatus string

const (
	NudgeStatusUnread    NudgeStatus = "unread"
	NudgeStatusRead      NudgeStatus = "read"
	NudgeStatusDismissed NudgeStatus = "dismissed"
	NudgeStatusActed     NudgeStatus = "acted"
)

type Nudge struct {
	BaseModel
	OrganizationID    uuid.UUID     `json:"organizationId" gorm:"not null"`
	Type              NudgeType     `json:"type"`
	Severity          NudgeSeverity `json:"severity"`
	Status            NudgeStatus   `json:"status" gorm:"default:'unread'"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	AIExplanation     string        `json:"aiExplanation" gorm:"column:ai_explanation"`
	SuggestedAction   string        `json:"suggestedAction"`
	ConfidenceScore   float64       `json:"confidenceScore"`
	CriticalityScore  int           `json:"criticalityScore"`
	ExpiresAt         *time.Time    `json:"expiresAt"`
	
	// Related entities
	RelatedProjectID *uuid.UUID `json:"relatedProjectId,omitempty"`
	RelatedTaskID    *uuid.UUID `json:"relatedTaskId,omitempty"`
	RelatedUserID    *uuid.UUID `json:"relatedUserId,omitempty"`
	
	// Metadata stored as JSONB
	Metrics JSONB `json:"metrics,omitempty" gorm:"type:jsonb"`
	
	// Relationships
	Organization  Organization    `json:"-" gorm:"foreignKey:OrganizationID"`
	RelatedProject *Project       `json:"relatedProject,omitempty" gorm:"foreignKey:RelatedProjectID"`
	RelatedTask    *Task          `json:"relatedTask,omitempty" gorm:"foreignKey:RelatedTaskID"`
	RelatedUser    *User          `json:"relatedUser,omitempty" gorm:"foreignKey:RelatedUserID"`
	Actions        []NudgeAction  `json:"actions,omitempty" gorm:"foreignKey:NudgeID"`
}

type NudgeAction struct {
	BaseModel
	NudgeID    uuid.UUID `json:"nudgeId" gorm:"not null"`
	UserID     uuid.UUID `json:"userId" gorm:"not null"`
	ActionType string    `json:"actionType"` // accept_suggestion, dismiss, custom_action, etc.
	Parameters JSONB    `json:"parameters,omitempty" gorm:"type:jsonb"`
	
	Nudge Nudge `json:"-" gorm:"foreignKey:NudgeID"`
	User  User  `json:"-" gorm:"foreignKey:UserID"`
}

// ===== Assignment Models =====

type AssignmentSuggestion struct {
	BaseModel
	TaskID           uuid.UUID `json:"taskId" gorm:"not null"`
	SuggestedUserID  uuid.UUID `json:"suggestedUserId" gorm:"not null"`
	TotalScore       int       `json:"totalScore"` // 0-100
	SkillMatchScore  int       `json:"skillMatchScore"` // 0-40
	AvailabilityScore int      `json:"availabilityScore"` // 0-30
	WorkloadScore    int       `json:"workloadScore"` // 0-20
	PerformanceScore int       `json:"performanceScore"` // 0-10
	Reasons          JSONB     `json:"reasons" gorm:"type:jsonb"`
	Warnings         JSONB     `json:"warnings,omitempty" gorm:"type:jsonb"`
	AIExplanation    string    `json:"aiExplanation"`
	Status           string    `json:"status" gorm:"default:'pending'"` // pending, accepted, rejected
	
	Task          Task `json:"-" gorm:"foreignKey:TaskID"`
	SuggestedUser User `json:"-" gorm:"foreignKey:SuggestedUserID"`
}

// ===== Workload Models =====

type WorkloadEntry struct {
	BaseModel
	OrganizationID       uuid.UUID `json:"organizationId" gorm:"not null"`
	UserID               uuid.UUID `json:"userId" gorm:"not null"`
	WeekStart            time.Time `json:"weekStart"`
	AllocationPercentage int       `json:"allocationPercentage"` // can exceed 100%
	AssignedTasks        int       `json:"assignedTasks"`
	TotalEstimatedHours  float64   `json:"totalEstimatedHours"`
	AvailableHours       float64   `json:"availableHours"`
	
	User User `json:"-" gorm:"foreignKey:UserID"`
}

// ===== Scenario Models =====

type ScenarioChangeType string

const (
	ScenarioChangeEmployeeLeave   ScenarioChangeType = "employee_leave"
	ScenarioChangeScopeChange     ScenarioChangeType = "scope_change"
	ScenarioChangeReallocation    ScenarioChangeType = "reallocation"
	ScenarioChangePriorityShift   ScenarioChangeType = "priority_shift"
)

type ScenarioStatus string

const (
	ScenarioStatusPending   ScenarioStatus = "pending"
	ScenarioStatusApproved  ScenarioStatus = "approved"
	ScenarioStatusRejected  ScenarioStatus = "rejected"
	ScenarioStatusModified  ScenarioStatus = "modified"
	ScenarioStatusApplied   ScenarioStatus = "applied"
)

type Scenario struct {
	BaseModel
	OrganizationID   uuid.UUID          `json:"organizationId" gorm:"not null"`
	Title            string             `json:"title"`
	Description      string             `json:"description"`
	ChangeType       ScenarioChangeType `json:"changeType"`
	Status           ScenarioStatus     `json:"status" gorm:"default:'pending'"`
	ProposedChanges  JSONB              `json:"proposedChanges" gorm:"type:jsonb"`
	CreatedByID      uuid.UUID          `json:"createdById"`
	DecidedByID      *uuid.UUID         `json:"decidedById,omitempty"`
	DecidedAt        *time.Time         `json:"decidedAt,omitempty"`
	
	// Relationships
	Organization    Organization         `json:"-" gorm:"foreignKey:OrganizationID"`
	CreatedBy       User                 `json:"-" gorm:"foreignKey:CreatedByID"`
	ImpactAnalysis  *ScenarioImpactAnalysis `json:"impactAnalysis,omitempty" gorm:"foreignKey:ScenarioID"`
}

type ScenarioImpactAnalysis struct {
	BaseModel
	ScenarioID         uuid.UUID `json:"scenarioId" gorm:"not null;uniqueIndex"`
	DelayHoursTotal    int       `json:"delayHoursTotal"`
	CostImpact         float64   `json:"costImpact"`
	AffectedProjectIDs JSONB     `json:"affectedProjectIds" gorm:"type:jsonb"`
	AffectedTaskIDs    JSONB     `json:"affectedTaskIds" gorm:"type:jsonb"`
	Recommendations    JSONB     `json:"recommendations" gorm:"type:jsonb"`
	TimelineComparison JSONB     `json:"timelineComparison" gorm:"type:jsonb"`
	
	Scenario Scenario `json:"-" gorm:"foreignKey:ScenarioID"`
}

// ===== JSONB Type Helper =====

type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (interface{}, error) {
	if j == nil {
		return nil, nil
	}
	return j, nil
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	// GORM/PostgreSQL will handle JSONB automatically
	return nil
}
