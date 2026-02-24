package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/SimpleAjax/Xephyr/internal/models"
)

// TaskRepository defines task data access operations
type TaskRepository interface {
	// Create creates a new task
	Create(ctx context.Context, task *models.Task) error

	// GetByID retrieves a task by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error)

	// Update updates a task
	Update(ctx context.Context, task *models.Task) error

	// Delete soft-deletes a task
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves tasks with pagination
	List(ctx context.Context, params ListParams) ([]models.Task, int64, error)

	// ListByProject retrieves tasks in a project
	ListByProject(ctx context.Context, projectID uuid.UUID, params ListParams) ([]models.Task, int64, error)

	// ListByAssignee retrieves tasks assigned to a user
	ListByAssignee(ctx context.Context, assigneeID uuid.UUID, params ListParams) ([]models.Task, int64, error)

	// ListByStatus retrieves tasks by status
	ListByStatus(ctx context.Context, status models.TaskStatus, params ListParams) ([]models.Task, int64, error)

	// ListSubtasks retrieves subtasks of a task
	ListSubtasks(ctx context.Context, parentTaskID uuid.UUID) ([]models.Task, error)

	// UpdateStatus updates task status
	UpdateStatus(ctx context.Context, taskID uuid.UUID, status models.TaskStatus) error

	// UpdateAssignee updates task assignee
	UpdateAssignee(ctx context.Context, taskID uuid.UUID, assigneeID *uuid.UUID) error

	// UpdateProgress updates task progress
	UpdateProgress(ctx context.Context, taskID uuid.UUID, progress int, actualHours float64) error

	// UpdatePriorityScore updates the calculated priority score
	UpdatePriorityScore(ctx context.Context, taskID uuid.UUID, score int) error

	// MarkAsCompleted marks a task as completed
	MarkAsCompleted(ctx context.Context, taskID uuid.UUID) error

	// ListByProjectAndStatus retrieves tasks in a project by status
	ListByProjectAndStatus(ctx context.Context, projectID uuid.UUID, status models.TaskStatus) ([]models.Task, error)

	// CountByProject counts tasks in a project
	CountByProject(ctx context.Context, projectID uuid.UUID) (int64, error)

	// CountByProjectAndStatus counts tasks by project and status
	CountByProjectAndStatus(ctx context.Context, projectID uuid.UUID, status models.TaskStatus) (int64, error)

	// GetUnassigned retrieves unassigned tasks in a project
	GetUnassigned(ctx context.Context, projectID uuid.UUID) ([]models.Task, error)

	// GetOverdue retrieves overdue tasks
	GetOverdue(ctx context.Context, orgID uuid.UUID) ([]models.Task, error)
}

// taskRepository implements TaskRepository
type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	var task models.Task
	if err := r.db.WithContext(ctx).
		Preload("Assignee").
		Preload("Skills.Skill").
		Preload("Subtasks").
		First(&task, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("task not found: %w", err)
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Task{}, "id = ?", id).Error
}

func (r *taskRepository) List(ctx context.Context, params ListParams) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Task{})
	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) ListByProject(ctx context.Context, projectID uuid.UUID, params ListParams) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Task{}).Where("project_id = ?", projectID)
	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) ListByAssignee(ctx context.Context, assigneeID uuid.UUID, params ListParams) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Task{}).Where("assignee_id = ?", assigneeID)
	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) ListByStatus(ctx context.Context, status models.TaskStatus, params ListParams) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Task{}).Where("status = ?", status)
	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) ListSubtasks(ctx context.Context, parentTaskID uuid.UUID) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("parent_task_id = ?", parentTaskID).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateStatus(ctx context.Context, taskID uuid.UUID, status models.TaskStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	if status == models.TaskStatusDone {
		now := time.Now()
		updates["completed_at"] = &now
		updates["progress"] = 100
	}
	
	return r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("id = ?", taskID).
		Updates(updates).Error
}

func (r *taskRepository) UpdateAssignee(ctx context.Context, taskID uuid.UUID, assigneeID *uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("id = ?", taskID).
		Update("assignee_id", assigneeID).Error
}

func (r *taskRepository) UpdateProgress(ctx context.Context, taskID uuid.UUID, progress int, actualHours float64) error {
	return r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"progress_percentage": progress,
			"actual_hours":        actualHours,
		}).Error
}

func (r *taskRepository) UpdatePriorityScore(ctx context.Context, taskID uuid.UUID, score int) error {
	return r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("id = ?", taskID).
		Update("priority_score", score).Error
}

func (r *taskRepository) MarkAsCompleted(ctx context.Context, taskID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"status":       models.TaskStatusDone,
			"completed_at": &now,
			"progress":     100,
		}).Error
}

func (r *taskRepository) ListByProjectAndStatus(ctx context.Context, projectID uuid.UUID, status models.TaskStatus) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND status = ?", projectID, status).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) CountByProject(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("project_id = ?", projectID).
		Count(&count).Error
	return count, err
}

func (r *taskRepository) CountByProjectAndStatus(ctx context.Context, projectID uuid.UUID, status models.TaskStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("project_id = ? AND status = ?", projectID, status).
		Count(&count).Error
	return count, err
}

func (r *taskRepository) GetUnassigned(ctx context.Context, projectID uuid.UUID) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND assignee_id IS NULL", projectID).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetOverdue(ctx context.Context, orgID uuid.UUID) ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()
	err := r.db.WithContext(ctx).
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("projects.organization_id = ?", orgID).
		Where("tasks.due_date < ? AND tasks.status != ?", now, models.TaskStatusDone).
		Find(&tasks).Error
	return tasks, err
}
