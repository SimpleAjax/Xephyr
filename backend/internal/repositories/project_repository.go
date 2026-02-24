package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/SimpleAjax/Xephyr/internal/models"
)

// ProjectRepository defines project data access operations
type ProjectRepository interface {
	// Create creates a new project
	Create(ctx context.Context, project *models.Project) error

	// GetByID retrieves a project by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error)

	// Update updates a project
	Update(ctx context.Context, project *models.Project) error

	// Delete soft-deletes a project
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves projects with pagination
	List(ctx context.Context, params ListParams) ([]models.Project, int64, error)

	// ListByOrganization retrieves projects in an organization
	ListByOrganization(ctx context.Context, orgID uuid.UUID, params ListParams) ([]models.Project, int64, error)

	// ListByStatus retrieves projects by status
	ListByStatus(ctx context.Context, status models.ProjectStatus, params ListParams) ([]models.Project, int64, error)

	// UpdateHealthScore updates project health score
	UpdateHealthScore(ctx context.Context, projectID uuid.UUID, score int) error

	// UpdateProgress updates project progress percentage
	UpdateProgress(ctx context.Context, projectID uuid.UUID, progress int) error

	// UpdatePriority updates project priority
	UpdatePriority(ctx context.Context, projectID uuid.UUID, priority int) error

	// AddMember adds a user to a project
	AddMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, role string) error

	// RemoveMember removes a user from a project
	RemoveMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error

	// IsMember checks if a user is a member of a project
	IsMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (bool, error)

	// CountTasks counts tasks in a project
	CountTasks(ctx context.Context, projectID uuid.UUID) (int64, error)
}

// projectRepository implements ProjectRepository
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).
		Preload("Members.User").
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_task_id IS NULL") // Only root tasks
		}).
		First(&project, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("project not found: %w", err)
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", id).Error
}

func (r *projectRepository) List(ctx context.Context, params ListParams) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Project{})
	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *projectRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, params ListParams) ([]models.Project, int64, error) {
	var projects []models.Project

	log.Printf("[ProjectRepository] Querying projects for organization_id: %s", orgID.String())

	// Try raw SQL to debug
	var rawCount int64
	r.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM projects WHERE organization_id = ?", orgID).Scan(&rawCount)
	log.Printf("[ProjectRepository] Raw SQL count: %d", rawCount)

	// Use GORM query
	scopedQuery := r.db.WithContext(ctx).Model(&models.Project{}).Where("organization_id = ?::uuid", orgID)
	var scopedTotal int64
	scopedQuery.Count(&scopedTotal)
	log.Printf("[ProjectRepository] GORM count: %d", scopedTotal)

	err := scopedQuery.Scopes(Paginate(params), Sort(params)).Find(&projects).Error
	if err != nil {
		log.Printf("[ProjectRepository] Error querying projects: %v", err)
		return nil, 0, err
	}

	log.Printf("[ProjectRepository] Found %d projects (slice len)", len(projects))
	for i, p := range projects {
		log.Printf("[ProjectRepository] Project %d: ID=%s, Name=%s", i, p.ID.String(), p.Name)
	}
	return projects, scopedTotal, nil
}

func (r *projectRepository) ListByStatus(ctx context.Context, status models.ProjectStatus, params ListParams) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Project{}).Where("status = ?", status)
	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *projectRepository) UpdateHealthScore(ctx context.Context, projectID uuid.UUID, score int) error {
	return r.db.WithContext(ctx).
		Model(&models.Project{}).
		Where("id = ?", projectID).
		Update("health_score", score).Error
}

func (r *projectRepository) UpdateProgress(ctx context.Context, projectID uuid.UUID, progress int) error {
	return r.db.WithContext(ctx).
		Model(&models.Project{}).
		Where("id = ?", projectID).
		Update("progress", progress).Error
}

func (r *projectRepository) UpdatePriority(ctx context.Context, projectID uuid.UUID, priority int) error {
	return r.db.WithContext(ctx).
		Model(&models.Project{}).
		Where("id = ?", projectID).
		Update("priority", priority).Error
}

func (r *projectRepository) AddMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, role string) error {
	member := models.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}
	return r.db.WithContext(ctx).Create(&member).Error
}

func (r *projectRepository) RemoveMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&models.ProjectMember{}).Error
}

func (r *projectRepository) IsMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *projectRepository) CountTasks(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("project_id = ?", projectID).
		Count(&count).Error
	return count, err
}
