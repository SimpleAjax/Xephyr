package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/SimpleAjax/Xephyr/internal/models"
)

// WorkloadRepository defines workload data access operations
type WorkloadRepository interface {
	// CreateOrUpdate creates or updates a workload entry
	CreateOrUpdate(ctx context.Context, entry *models.WorkloadEntry) error

	// GetByUserAndWeek retrieves workload for a user in a specific week
	GetByUserAndWeek(ctx context.Context, userID uuid.UUID, weekStart time.Time) (*models.WorkloadEntry, error)

	// ListByOrganization retrieves workload entries for an organization
	ListByOrganization(ctx context.Context, orgID uuid.UUID, weekStart time.Time) ([]models.WorkloadEntry, error)

	// ListByUser retrieves workload history for a user
	ListByUser(ctx context.Context, userID uuid.UUID, fromDate, toDate time.Time) ([]models.WorkloadEntry, error)

	// GetTeamWorkload retrieves aggregated team workload
	GetTeamWorkload(ctx context.Context, orgID uuid.UUID, weekStart time.Time) (*TeamWorkload, error)

	// UpdateAllocation updates user's allocation percentage
	UpdateAllocation(ctx context.Context, userID uuid.UUID, weekStart time.Time, allocation int) error

	// DeleteOldEntries removes workload entries older than threshold
	DeleteOldEntries(ctx context.Context, olderThan time.Duration) error
}

// TeamWorkload represents aggregated team workload
type TeamWorkload struct {
	WeekStarting    time.Time              `json:"weekStarting"`
	TeamCapacity    float64                `json:"teamCapacity"`
	TeamAllocation  float64                `json:"teamAllocation"`
	UtilizationRate float64                `json:"utilizationRate"`
	MemberWorkloads []MemberWorkloadDetail `json:"memberWorkloads"`
}

// MemberWorkloadDetail represents detailed workload for a member
type MemberWorkloadDetail struct {
	UserID              uuid.UUID `json:"userId"`
	UserName            string    `json:"userName"`
	AllocationPercentage int      `json:"allocationPercentage"`
	AssignedHours       float64   `json:"assignedHours"`
	CapacityHours       float64   `json:"capacityHours"`
	AssignedTasks       int       `json:"assignedTasks"`
}

// workloadRepository implements WorkloadRepository
type workloadRepository struct {
	db *gorm.DB
}

// NewWorkloadRepository creates a new workload repository
func NewWorkloadRepository(db *gorm.DB) WorkloadRepository {
	return &workloadRepository{db: db}
}

func (r *workloadRepository) CreateOrUpdate(ctx context.Context, entry *models.WorkloadEntry) error {
	// Try to update existing entry
	result := r.db.WithContext(ctx).
		Model(&models.WorkloadEntry{}).
		Where("user_id = ? AND week_start = ?", entry.UserID, entry.WeekStart).
		Updates(map[string]interface{}{
			"allocation_percentage": entry.AllocationPercentage,
			"assigned_tasks":        entry.AssignedTasks,
			"total_estimated_hours": entry.TotalEstimatedHours,
			"available_hours":       entry.AvailableHours,
		})

	if result.Error != nil {
		return result.Error
	}

	// If no rows updated, create new
	if result.RowsAffected == 0 {
		return r.db.WithContext(ctx).Create(entry).Error
	}

	return nil
}

func (r *workloadRepository) GetByUserAndWeek(ctx context.Context, userID uuid.UUID, weekStart time.Time) (*models.WorkloadEntry, error) {
	var entry models.WorkloadEntry
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ? AND week_start = ?", userID, weekStart).
		First(&entry).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *workloadRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, weekStart time.Time) ([]models.WorkloadEntry, error) {
	var entries []models.WorkloadEntry
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("organization_id = ? AND week_start = ?", orgID, weekStart).
		Find(&entries).Error
	return entries, err
}

func (r *workloadRepository) ListByUser(ctx context.Context, userID uuid.UUID, fromDate, toDate time.Time) ([]models.WorkloadEntry, error) {
	var entries []models.WorkloadEntry
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND week_start BETWEEN ? AND ?", userID, fromDate, toDate).
		Order("week_start ASC").
		Find(&entries).Error
	return entries, err
}

func (r *workloadRepository) GetTeamWorkload(ctx context.Context, orgID uuid.UUID, weekStart time.Time) (*TeamWorkload, error) {
	var entries []models.WorkloadEntry
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("organization_id = ? AND week_start = ?", orgID, weekStart).
		Find(&entries).Error
	if err != nil {
		return nil, err
	}

	tw := &TeamWorkload{
		WeekStarting:    weekStart,
		MemberWorkloads: make([]MemberWorkloadDetail, 0, len(entries)),
	}

	for _, entry := range entries {
		tw.TeamCapacity += entry.AvailableHours
		tw.TeamAllocation += entry.TotalEstimatedHours
		
		detail := MemberWorkloadDetail{
			UserID:               entry.UserID,
			UserName:             entry.User.Name,
			AllocationPercentage: entry.AllocationPercentage,
			AssignedHours:        entry.TotalEstimatedHours,
			CapacityHours:        entry.AvailableHours,
			AssignedTasks:        entry.AssignedTasks,
		}
		
		tw.MemberWorkloads = append(tw.MemberWorkloads, detail)
	}

	if tw.TeamCapacity > 0 {
		tw.UtilizationRate = tw.TeamAllocation / tw.TeamCapacity
	}

	return tw, nil
}

func (r *workloadRepository) UpdateAllocation(ctx context.Context, userID uuid.UUID, weekStart time.Time, allocation int) error {
	return r.db.WithContext(ctx).
		Model(&models.WorkloadEntry{}).
		Where("user_id = ? AND week_start = ?", userID, weekStart).
		Update("allocation_percentage", allocation).Error
}

func (r *workloadRepository) DeleteOldEntries(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	return r.db.WithContext(ctx).
		Where("week_start < ?", cutoff).
		Delete(&models.WorkloadEntry{}).Error
}
