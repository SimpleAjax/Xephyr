package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/SimpleAjax/Xephyr/internal/models"
)

// UserRepository defines user data access operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// Update updates a user
	Update(ctx context.Context, user *models.User) error

	// Delete soft-deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves users with pagination
	List(ctx context.Context, params ListParams) ([]models.User, int64, error)

	// ListByOrganization retrieves users in an organization
	ListByOrganization(ctx context.Context, orgID uuid.UUID, params ListParams) ([]models.User, int64, error)

	// GetByOrganizationAndRole retrieves users by role in an organization
	GetByOrganizationAndRole(ctx context.Context, orgID uuid.UUID, role models.UserRole) ([]models.User, error)

	// UpdatePassword updates user password
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error

	// Exists checks if a user exists
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Skills.Skill").First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

func (r *userRepository) List(ctx context.Context, params ListParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})
	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, params ListParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.User{}).
		Joins("JOIN organization_members ON users.id = organization_members.user_id").
		Where("organization_members.organization_id = ?", orgID)

	query.Count(&total)

	err := query.Scopes(Paginate(params), Sort(params)).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetByOrganizationAndRole(ctx context.Context, orgID uuid.UUID, role models.UserRole) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Joins("JOIN organization_members ON users.id = organization_members.user_id").
		Where("organization_members.organization_id = ? AND organization_members.role = ?", orgID, role).
		Find(&users).Error
	return users, err
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("password_hash", passwordHash).Error
}

func (r *userRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
