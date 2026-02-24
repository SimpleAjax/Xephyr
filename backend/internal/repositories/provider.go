package repositories

import (
	"gorm.io/gorm"
)

// Provider holds all repository instances
type Provider struct {
	User         UserRepository
	Organization OrganizationRepository
	Project      ProjectRepository
	Task         TaskRepository
	Nudge        NudgeRepository
	Assignment   AssignmentRepository
	Workload     WorkloadRepository
	Scenario     ScenarioRepository
	Dependency   DependencyRepository
}

// NewProvider creates a new repository provider with all repositories
func NewProvider(db *gorm.DB) *Provider {
	return &Provider{
		User:         NewUserRepository(db),
		Organization: NewOrganizationRepository(db),
		Project:      NewProjectRepository(db),
		Task:         NewTaskRepository(db),
		Nudge:        NewNudgeRepository(db),
		Assignment:   NewAssignmentRepository(db),
		Workload:     NewWorkloadRepository(db),
		Scenario:     NewScenarioRepository(db),
		Dependency:   NewDependencyRepository(db),
	}
}

// Repositories interface for easy mocking in tests
type Repositories interface {
	GetUser() UserRepository
	GetOrganization() OrganizationRepository
	GetProject() ProjectRepository
	GetTask() TaskRepository
	GetNudge() NudgeRepository
	GetAssignment() AssignmentRepository
	GetWorkload() WorkloadRepository
	GetScenario() ScenarioRepository
	GetDependency() DependencyRepository
}

// Ensure Provider implements Repositories
var _ Repositories = (*Provider)(nil)

// GetUser returns the user repository
func (p *Provider) GetUser() UserRepository {
	return p.User
}

// GetOrganization returns the organization repository
func (p *Provider) GetOrganization() OrganizationRepository {
	return p.Organization
}

// GetProject returns the project repository
func (p *Provider) GetProject() ProjectRepository {
	return p.Project
}

// GetTask returns the task repository
func (p *Provider) GetTask() TaskRepository {
	return p.Task
}

// GetNudge returns the nudge repository
func (p *Provider) GetNudge() NudgeRepository {
	return p.Nudge
}

// GetAssignment returns the assignment repository
func (p *Provider) GetAssignment() AssignmentRepository {
	return p.Assignment
}

// GetWorkload returns the workload repository
func (p *Provider) GetWorkload() WorkloadRepository {
	return p.Workload
}

// GetScenario returns the scenario repository
func (p *Provider) GetScenario() ScenarioRepository {
	return p.Scenario
}

// GetDependency returns the dependency repository
func (p *Provider) GetDependency() DependencyRepository {
	return p.Dependency
}
