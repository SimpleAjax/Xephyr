package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// DefaultDBConfig returns default configuration
type DefaultDBConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	Database string `env:"DB_NAME" envDefault:"xephyr"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

// Repository provides database connection and transaction support
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository with database connection
func NewRepository(config DBConfig) (*Repository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return &Repository{db: db}, nil
}

// DB returns the underlying gorm.DB
func (r *Repository) DB() *gorm.DB {
	return r.db
}

// WithTransaction executes a function within a database transaction
func (r *Repository) WithTransaction(ctx context.Context, fn func(*Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&Repository{db: tx})
	})
}

// WithContext returns a new repository with context
func (r *Repository) WithContext(ctx context.Context) *Repository {
	return &Repository{db: r.db.WithContext(ctx)}
}

// Migrate runs database migrations
func (r *Repository) Migrate(models ...interface{}) error {
	return r.db.AutoMigrate(models...)
}

// Health checks database connectivity
func (r *Repository) Health() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close closes the database connection
func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// ListParams provides common pagination and sorting parameters
type ListParams struct {
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string
}

// DefaultListParams returns default list parameters
func DefaultListParams() ListParams {
	return ListParams{
		Offset:    0,
		Limit:     20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// Paginate applies pagination to a query
func Paginate(params ListParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		limit := params.Limit
		if limit <= 0 {
			limit = 20
		}
		if limit > 100 {
			limit = 100
		}
		return db.Offset(params.Offset).Limit(limit)
	}
}

// Sort applies sorting to a query
func Sort(params ListParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		sortBy := params.SortBy
		if sortBy == "" {
			sortBy = "created_at"
		}
		sortOrder := params.SortOrder
		if sortOrder == "" {
			sortOrder = "desc"
		}
		order := sortBy + " " + sortOrder
		return db.Order(order)
	}
}

// logQuery logs slow queries (optional helper)
func logQuery(start time.Time, query string, args ...interface{}) {
	duration := time.Since(start)
	if duration > 100*time.Millisecond {
		log.Printf("Slow query (%v): %s", duration, query)
	}
}
