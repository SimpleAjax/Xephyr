package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/SimpleAjax/Xephyr/internal/models"
)

// OrganizationRepository defines organization data access operations
type OrganizationRepository interface {
	// Create creates a new organization
	Create(ctx context.Context, org *models.Organization) error

	// GetByID retrieves an organization by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)

	// GetBySlug retrieves an organization by slug
	GetBySlug(ctx context.Context, slug string) (*models.Organization, error)

	// Update updates an organization
	Update(ctx context.Context, org *models.Organization) error

	// Delete soft-deletes an organization
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves organizations with pagination
	List(ctx context.Context, params ListParams) ([]models.Organization, int64, error)

	// AddMember adds a user to an organization
	AddMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role models.UserRole) error

	// RemoveMember removes a user from an organization
	RemoveMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error

	// UpdateMemberRole updates a member's role
	UpdateMemberRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role models.UserRole) error

	// IsMember checks if a user is a member of an organization
	IsMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (bool, error)

	// GetMemberRole gets a user's role in an organization
	GetMemberRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (models.UserRole, error)
}

// organizationRepository implements OrganizationRepository
type organizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, org *models.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).
		Preload("Members.User").
		Preload("Projects").
		First(&org, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("organization not found: %w", err)
		}
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) GetBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).First(&org, "slug = ?", slug).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("organization not found: %w", err)
		}
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) Update(ctx context.Context, org *models.Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Organization{}, "id = ?", id).Error
}

func (r *organizationRepository) List(ctx context.Context, params ListParams) ([]models.Organization, int64, error) {
	var orgs []models.Organization
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Organization{})
	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&orgs).Error
	if err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

func (r *organizationRepository) AddMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role models.UserRole) error {
	member := models.OrganizationMember{
		OrganizationID: orgID,
		UserID:         userID,
		Role:           role,
	}
	return r.db.WithContext(ctx).Create(&member).Error
}

func (r *organizationRepository) RemoveMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Delete(&models.OrganizationMember{}).Error
}

func (r *organizationRepository) UpdateMemberRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role models.UserRole) error {
	return r.db.WithContext(ctx).
		Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Update("role", role).Error
}

func (r *organizationRepository) IsMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *organizationRepository) GetMemberRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (models.UserRole, error) {
	var member models.OrganizationMember
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		First(&member).Error
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
