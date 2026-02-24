package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/SimpleAjax/Xephyr/internal/models"
)

// ScenarioRepository defines scenario data access operations
type ScenarioRepository interface {
	// Create creates a new scenario
	Create(ctx context.Context, scenario *models.Scenario) error

	// GetByID retrieves a scenario by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Scenario, error)

	// Update updates a scenario
	Update(ctx context.Context, scenario *models.Scenario) error

	// Delete soft-deletes a scenario
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves scenarios with filters
	List(ctx context.Context, orgID uuid.UUID, filters ScenarioFilters, params ListParams) ([]models.Scenario, int64, error)

	// UpdateStatus updates scenario status
	UpdateStatus(ctx context.Context, scenarioID uuid.UUID, status models.ScenarioStatus) error

	// CreateImpactAnalysis creates impact analysis for a scenario
	CreateImpactAnalysis(ctx context.Context, analysis *models.ScenarioImpactAnalysis) error

	// GetImpactAnalysis retrieves impact analysis for a scenario
	GetImpactAnalysis(ctx context.Context, scenarioID uuid.UUID) (*models.ScenarioImpactAnalysis, error)

	// UpdateImpactAnalysis updates impact analysis
	UpdateImpactAnalysis(ctx context.Context, analysis *models.ScenarioImpactAnalysis) error

	// GetByStatus retrieves scenarios by status
	GetByStatus(ctx context.Context, orgID uuid.UUID, status models.ScenarioStatus) ([]models.Scenario, error)
}

// ScenarioFilters provides filtering options for scenarios
type ScenarioFilters struct {
	Status    *models.ScenarioStatus
	ChangeType *models.ScenarioChangeType
	CreatedBy *uuid.UUID
}

// scenarioRepository implements ScenarioRepository
type scenarioRepository struct {
	db *gorm.DB
}

// NewScenarioRepository creates a new scenario repository
func NewScenarioRepository(db *gorm.DB) ScenarioRepository {
	return &scenarioRepository{db: db}
}

func (r *scenarioRepository) Create(ctx context.Context, scenario *models.Scenario) error {
	return r.db.WithContext(ctx).Create(scenario).Error
}

func (r *scenarioRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Scenario, error) {
	var scenario models.Scenario
	if err := r.db.WithContext(ctx).
		Preload("ImpactAnalysis").
		Preload("CreatedBy").
		First(&scenario, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("scenario not found: %w", err)
		}
		return nil, err
	}
	return &scenario, nil
}

func (r *scenarioRepository) Update(ctx context.Context, scenario *models.Scenario) error {
	return r.db.WithContext(ctx).Save(scenario).Error
}

func (r *scenarioRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Scenario{}, "id = ?", id).Error
}

func (r *scenarioRepository) List(ctx context.Context, orgID uuid.UUID, filters ScenarioFilters, params ListParams) ([]models.Scenario, int64, error) {
	var scenarios []models.Scenario
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Scenario{}).Where("organization_id = ?", orgID)

	// Apply filters
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.ChangeType != nil {
		query = query.Where("change_type = ?", *filters.ChangeType)
	}
	if filters.CreatedBy != nil {
		query = query.Where("created_by_id = ?", *filters.CreatedBy)
	}

	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&scenarios).Error
	if err != nil {
		return nil, 0, err
	}

	return scenarios, total, nil
}

func (r *scenarioRepository) UpdateStatus(ctx context.Context, scenarioID uuid.UUID, status models.ScenarioStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Scenario{}).
		Where("id = ?", scenarioID).
		Update("status", status).Error
}

func (r *scenarioRepository) CreateImpactAnalysis(ctx context.Context, analysis *models.ScenarioImpactAnalysis) error {
	return r.db.WithContext(ctx).Create(analysis).Error
}

func (r *scenarioRepository) GetImpactAnalysis(ctx context.Context, scenarioID uuid.UUID) (*models.ScenarioImpactAnalysis, error) {
	var analysis models.ScenarioImpactAnalysis
	if err := r.db.WithContext(ctx).
		Where("scenario_id = ?", scenarioID).
		First(&analysis).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("impact analysis not found: %w", err)
		}
		return nil, err
	}
	return &analysis, nil
}

func (r *scenarioRepository) UpdateImpactAnalysis(ctx context.Context, analysis *models.ScenarioImpactAnalysis) error {
	return r.db.WithContext(ctx).Save(analysis).Error
}

func (r *scenarioRepository) GetByStatus(ctx context.Context, orgID uuid.UUID, status models.ScenarioStatus) ([]models.Scenario, error) {
	var scenarios []models.Scenario
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND status = ?", orgID, status).
		Find(&scenarios).Error
	return scenarios, err
}
