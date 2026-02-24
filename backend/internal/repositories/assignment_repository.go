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

// AssignmentRepository defines assignment-related data access operations
type AssignmentRepository interface {
	// CreateSuggestion creates a new assignment suggestion
	CreateSuggestion(ctx context.Context, suggestion *models.AssignmentSuggestion) error

	// GetSuggestionByID retrieves a suggestion by ID
	GetSuggestionByID(ctx context.Context, id uuid.UUID) (*models.AssignmentSuggestion, error)

	// GetSuggestionsByTask retrieves suggestions for a task
	GetSuggestionsByTask(ctx context.Context, taskID uuid.UUID) ([]models.AssignmentSuggestion, error)

	// UpdateSuggestionStatus updates suggestion status
	UpdateSuggestionStatus(ctx context.Context, suggestionID uuid.UUID, status string) error

	// DeleteOldSuggestions removes old pending suggestions
	DeleteOldSuggestions(ctx context.Context, olderThan time.Duration) error

	// CreateAssignmentHistory logs an assignment
	CreateAssignmentHistory(ctx context.Context, taskID uuid.UUID, fromUserID, toUserID *uuid.UUID, assignedBy uuid.UUID) error

	// GetAssignmentHistory retrieves assignment history for a task
	GetAssignmentHistory(ctx context.Context, taskID uuid.UUID) ([]AssignmentHistory, error)
}

// AssignmentHistory represents an assignment history entry
type AssignmentHistory struct {
	ID          uuid.UUID  `json:"id"`
	TaskID      uuid.UUID  `json:"taskId"`
	FromUserID  *uuid.UUID `json:"fromUserId,omitempty"`
	ToUserID    *uuid.UUID `json:"toUserId,omitempty"`
	AssignedBy  uuid.UUID  `json:"assignedBy"`
	AssignedAt  time.Time  `json:"assignedAt"`
}

// assignmentRepository implements AssignmentRepository
type assignmentRepository struct {
	db *gorm.DB
}

// NewAssignmentRepository creates a new assignment repository
func NewAssignmentRepository(db *gorm.DB) AssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) CreateSuggestion(ctx context.Context, suggestion *models.AssignmentSuggestion) error {
	return r.db.WithContext(ctx).Create(suggestion).Error
}

func (r *assignmentRepository) GetSuggestionByID(ctx context.Context, id uuid.UUID) (*models.AssignmentSuggestion, error) {
	var suggestion models.AssignmentSuggestion
	if err := r.db.WithContext(ctx).
		Preload("SuggestedUser").
		Preload("Task").
		First(&suggestion, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("suggestion not found: %w", err)
		}
		return nil, err
	}
	return &suggestion, nil
}

func (r *assignmentRepository) GetSuggestionsByTask(ctx context.Context, taskID uuid.UUID) ([]models.AssignmentSuggestion, error) {
	var suggestions []models.AssignmentSuggestion
	err := r.db.WithContext(ctx).
		Preload("SuggestedUser").
		Where("task_id = ?", taskID).
		Order("total_score DESC").
		Find(&suggestions).Error
	return suggestions, err
}

func (r *assignmentRepository) UpdateSuggestionStatus(ctx context.Context, suggestionID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.AssignmentSuggestion{}).
		Where("id = ?", suggestionID).
		Update("status", status).Error
}

func (r *assignmentRepository) DeleteOldSuggestions(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	return r.db.WithContext(ctx).
		Where("status = ? AND created_at < ?", "pending", cutoff).
		Delete(&models.AssignmentSuggestion{}).Error
}

func (r *assignmentRepository) CreateAssignmentHistory(ctx context.Context, taskID uuid.UUID, fromUserID, toUserID *uuid.UUID, assignedBy uuid.UUID) error {
	// For now, we'll store this in a JSON field or separate table
	// This is a placeholder implementation
	return nil
}

func (r *assignmentRepository) GetAssignmentHistory(ctx context.Context, taskID uuid.UUID) ([]AssignmentHistory, error) {
	// Placeholder - would query assignment history table
	return []AssignmentHistory{}, nil
}
