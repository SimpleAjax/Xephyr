package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/SimpleAjax/Xephyr/internal/models"
)

// DependencyRepository defines dependency data access operations
type DependencyRepository interface {
	// Create creates a new dependency
	Create(ctx context.Context, dep *models.TaskDependency) error

	// GetByID retrieves a dependency by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.TaskDependency, error)

	// Delete removes a dependency
	Delete(ctx context.Context, id uuid.UUID) error

	// ListByTask retrieves dependencies for a task
	ListByTask(ctx context.Context, taskID uuid.UUID) ([]models.TaskDependency, error)

	// ListByDependsOn retrieves tasks that depend on a given task
	ListByDependsOn(ctx context.Context, dependsOnTaskID uuid.UUID) ([]models.TaskDependency, error)

	// ListByProject retrieves all dependencies in a project
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.TaskDependency, error)

	// HasDependency checks if a dependency exists
	HasDependency(ctx context.Context, taskID, dependsOnTaskID uuid.UUID) (bool, error)

	// WouldCreateCycle checks if adding a dependency would create a cycle
	WouldCreateCycle(ctx context.Context, taskID, dependsOnTaskID uuid.UUID) (bool, error)

	// GetDependentCount counts how many tasks depend on this task
	GetDependentCount(ctx context.Context, taskID uuid.UUID) (int64, error)

	// GetDependencyCount counts how many tasks this task depends on
	GetDependencyCount(ctx context.Context, taskID uuid.UUID) (int64, error)

	// DeleteByTask removes all dependencies for a task
	DeleteByTask(ctx context.Context, taskID uuid.UUID) error

	// UpdateLag updates the lag hours for a dependency
	UpdateLag(ctx context.Context, dependencyID uuid.UUID, lagHours int) error
}

// dependencyRepository implements DependencyRepository
type dependencyRepository struct {
	db *gorm.DB
}

// NewDependencyRepository creates a new dependency repository
func NewDependencyRepository(db *gorm.DB) DependencyRepository {
	return &dependencyRepository{db: db}
}

func (r *dependencyRepository) Create(ctx context.Context, dep *models.TaskDependency) error {
	return r.db.WithContext(ctx).Create(dep).Error
}

func (r *dependencyRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TaskDependency, error) {
	var dep models.TaskDependency
	if err := r.db.WithContext(ctx).
		Preload("Task").
		Preload("DependsOnTask").
		First(&dep, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("dependency not found: %w", err)
		}
		return nil, err
	}
	return &dep, nil
}

func (r *dependencyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.TaskDependency{}, "id = ?", id).Error
}

func (r *dependencyRepository) ListByTask(ctx context.Context, taskID uuid.UUID) ([]models.TaskDependency, error) {
	var deps []models.TaskDependency
	err := r.db.WithContext(ctx).
		Preload("DependsOnTask").
		Where("task_id = ?", taskID).
		Find(&deps).Error
	return deps, err
}

func (r *dependencyRepository) ListByDependsOn(ctx context.Context, dependsOnTaskID uuid.UUID) ([]models.TaskDependency, error) {
	var deps []models.TaskDependency
	err := r.db.WithContext(ctx).
		Preload("Task").
		Where("depends_on_task_id = ?", dependsOnTaskID).
		Find(&deps).Error
	return deps, err
}

func (r *dependencyRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.TaskDependency, error) {
	var deps []models.TaskDependency
	err := r.db.WithContext(ctx).
		Joins("JOIN tasks ON task_dependencies.task_id = tasks.id").
		Where("tasks.project_id = ?", projectID).
		Find(&deps).Error
	return deps, err
}

func (r *dependencyRepository) HasDependency(ctx context.Context, taskID, dependsOnTaskID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.TaskDependency{}).
		Where("task_id = ? AND depends_on_task_id = ?", taskID, dependsOnTaskID).
		Count(&count).Error
	return count > 0, err
}

// WouldCreateCycle checks if adding a dependency would create a cycle using CTE
func (r *dependencyRepository) WouldCreateCycle(ctx context.Context, taskID, dependsOnTaskID uuid.UUID) (bool, error) {
	// Direct cycle check
	if taskID == dependsOnTaskID {
		return true, nil
	}

	// Check if dependsOnTaskID already depends on taskID (would create cycle)
	var count int64
	// This is a simplified check - for a full check you'd need a recursive CTE
	err := r.db.WithContext(ctx).
		Raw(`
			WITH RECURSIVE dependency_chain AS (
				SELECT task_id, depends_on_task_id
				FROM task_dependencies
				WHERE task_id = ?
				UNION
				SELECT td.task_id, td.depends_on_task_id
				FROM task_dependencies td
				JOIN dependency_chain dc ON td.task_id = dc.depends_on_task_id
			)
			SELECT COUNT(*) FROM dependency_chain WHERE depends_on_task_id = ?
		`, dependsOnTaskID, taskID).
		Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *dependencyRepository) GetDependentCount(ctx context.Context, taskID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.TaskDependency{}).
		Where("depends_on_task_id = ?", taskID).
		Count(&count).Error
	return count, err
}

func (r *dependencyRepository) GetDependencyCount(ctx context.Context, taskID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.TaskDependency{}).
		Where("task_id = ?", taskID).
		Count(&count).Error
	return count, err
}

func (r *dependencyRepository) DeleteByTask(ctx context.Context, taskID uuid.UUID) error {
	// Delete where task is either the dependent or the dependency
	return r.db.WithContext(ctx).
		Where("task_id = ? OR depends_on_task_id = ?", taskID, taskID).
		Delete(&models.TaskDependency{}).Error
}

func (r *dependencyRepository) UpdateLag(ctx context.Context, dependencyID uuid.UUID, lagHours int) error {
	return r.db.WithContext(ctx).
		Model(&models.TaskDependency{}).
		Where("id = ?", dependencyID).
		Update("lag_hours", lagHours).Error
}
