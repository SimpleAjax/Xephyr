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

// NudgeRepository defines nudge data access operations
type NudgeRepository interface {
	// Create creates a new nudge
	Create(ctx context.Context, nudge *models.Nudge) error

	// GetByID retrieves a nudge by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Nudge, error)

	// Update updates a nudge
	Update(ctx context.Context, nudge *models.Nudge) error

	// Delete soft-deletes a nudge
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves nudges with filters
	List(ctx context.Context, orgID uuid.UUID, filters NudgeFilters, params ListParams) ([]models.Nudge, int64, error)

	// ListByUser retrieves nudges for a specific user
	ListByUser(ctx context.Context, userID uuid.UUID, params ListParams) ([]models.Nudge, int64, error)

	// UpdateStatus updates nudge status
	UpdateStatus(ctx context.Context, nudgeID uuid.UUID, status models.NudgeStatus) error

	// CountByStatus counts nudges by status
	CountByStatus(ctx context.Context, orgID uuid.UUID, status models.NudgeStatus) (int64, error)

	// GetStats retrieves nudge statistics
	GetStats(ctx context.Context, orgID uuid.UUID, period time.Duration) (*NudgeStats, error)

	// ExpireOldNudges marks expired nudges
	ExpireOldNudges(ctx context.Context) error

	// DeleteOldNudges permanently deletes old dismissed/acted nudges
	DeleteOldNudges(ctx context.Context, olderThan time.Duration) error
}

// NudgeFilters provides filtering options for nudges
type NudgeFilters struct {
	Status   *models.NudgeStatus
	Severity *models.NudgeSeverity
	Type     *models.NudgeType
	UserID   *uuid.UUID
	TaskID   *uuid.UUID
	ProjectID *uuid.UUID
}

// NudgeStats represents nudge statistics
type NudgeStats struct {
	Total      int64
	Unread     int64
	Read       int64
	Acted      int64
	Dismissed  int64
	BySeverity map[string]int64
	ByType     map[string]int64
}

// nudgeRepository implements NudgeRepository
type nudgeRepository struct {
	db *gorm.DB
}

// NewNudgeRepository creates a new nudge repository
func NewNudgeRepository(db *gorm.DB) NudgeRepository {
	return &nudgeRepository{db: db}
}

func (r *nudgeRepository) Create(ctx context.Context, nudge *models.Nudge) error {
	return r.db.WithContext(ctx).Create(nudge).Error
}

func (r *nudgeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Nudge, error) {
	var nudge models.Nudge
	if err := r.db.WithContext(ctx).
		Preload("RelatedUser").
		Preload("RelatedTask").
		Preload("RelatedProject").
		Preload("Actions.User").
		First(&nudge, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("nudge not found: %w", err)
		}
		return nil, err
	}
	return &nudge, nil
}

func (r *nudgeRepository) Update(ctx context.Context, nudge *models.Nudge) error {
	return r.db.WithContext(ctx).Save(nudge).Error
}

func (r *nudgeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Nudge{}, "id = ?", id).Error
}

func (r *nudgeRepository) List(ctx context.Context, orgID uuid.UUID, filters NudgeFilters, params ListParams) ([]models.Nudge, int64, error) {
	var nudges []models.Nudge
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Nudge{}).Where("organization_id = ?", orgID)

	// Apply filters
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.Severity != nil {
		query = query.Where("severity = ?", *filters.Severity)
	}
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.UserID != nil {
		query = query.Where("related_user_id = ?", *filters.UserID)
	}
	if filters.TaskID != nil {
		query = query.Where("related_task_id = ?", *filters.TaskID)
	}
	if filters.ProjectID != nil {
		query = query.Where("related_project_id = ?", *filters.ProjectID)
	}

	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&nudges).Error
	if err != nil {
		return nil, 0, err
	}

	return nudges, total, nil
}

func (r *nudgeRepository) ListByUser(ctx context.Context, userID uuid.UUID, params ListParams) ([]models.Nudge, int64, error) {
	var nudges []models.Nudge
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Nudge{}).Where("related_user_id = ?", userID)
	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&nudges).Error
	if err != nil {
		return nil, 0, err
	}

	return nudges, total, nil
}

func (r *nudgeRepository) UpdateStatus(ctx context.Context, nudgeID uuid.UUID, status models.NudgeStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Nudge{}).
		Where("id = ?", nudgeID).
		Update("status", status).Error
}

func (r *nudgeRepository) CountByStatus(ctx context.Context, orgID uuid.UUID, status models.NudgeStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Nudge{}).
		Where("organization_id = ? AND status = ?", orgID, status).
		Count(&count).Error
	return count, err
}

func (r *nudgeRepository) GetStats(ctx context.Context, orgID uuid.UUID, period time.Duration) (*NudgeStats, error) {
	stats := &NudgeStats{
		BySeverity: make(map[string]int64),
		ByType:     make(map[string]int64),
	}

	// Get total counts by status
	var results []struct {
		Status string
		Count  int64
	}
	
	since := time.Now().Add(-period)
	err := r.db.WithContext(ctx).
		Model(&models.Nudge{}).
		Select("status, COUNT(*) as count").
		Where("organization_id = ? AND created_at > ?", orgID, since).
		Group("status").
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}

	for _, r := range results {
		switch models.NudgeStatus(r.Status) {
		case models.NudgeStatusUnread:
			stats.Unread = r.Count
		case models.NudgeStatusRead:
			stats.Read = r.Count
		case models.NudgeStatusActed:
			stats.Acted = r.Count
		case models.NudgeStatusDismissed:
			stats.Dismissed = r.Count
		}
		stats.Total += r.Count
	}

	// Get counts by severity
	var severityResults []struct {
		Severity string
		Count    int64
	}
	
	err = r.db.WithContext(ctx).
		Model(&models.Nudge{}).
		Select("severity, COUNT(*) as count").
		Where("organization_id = ? AND created_at > ?", orgID, since).
		Group("severity").
		Scan(&severityResults).Error
	
	if err != nil {
		return nil, err
	}

	for _, r := range severityResults {
		stats.BySeverity[r.Severity] = r.Count
	}

	// Get counts by type
	var typeResults []struct {
		Type  string
		Count int64
	}
	
	err = r.db.WithContext(ctx).
		Model(&models.Nudge{}).
		Select("type, COUNT(*) as count").
		Where("organization_id = ? AND created_at > ?", orgID, since).
		Group("type").
		Scan(&typeResults).Error
	
	if err != nil {
		return nil, err
	}

	for _, r := range typeResults {
		stats.ByType[r.Type] = r.Count
	}

	return stats, nil
}

func (r *nudgeRepository) ExpireOldNudges(ctx context.Context) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Nudge{}).
		Where("expires_at < ? AND status = ?", now, models.NudgeStatusUnread).
		Update("status", models.NudgeStatusDismissed).Error
}

func (r *nudgeRepository) DeleteOldNudges(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	return r.db.WithContext(ctx).
		Where("(status = ? OR status = ?) AND updated_at < ?", 
			models.NudgeStatusDismissed, models.NudgeStatusActed, cutoff).
		Delete(&models.Nudge{}).Error
}
