package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/internal/repositories"
	"github.com/SimpleAjax/Xephyr/internal/routes"

	_ "github.com/SimpleAjax/Xephyr/docs"
)

// @title           Xephyr API
// @version         1.0
// @description     Xephyr Project Management Intelligence API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Get configuration from environment
	port := getEnv("PORT", "8080")
	// env := getEnv("ENV", "development")

	// Initialize database connection
	dbConfig := repositories.DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "xephyr"),
		Password: getEnv("DB_PASSWORD", "xephyr123"),
		Database: getEnv("DB_NAME", "xephyr"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	repo, err := repositories.NewRepository(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()

	// Run database migrations
	log.Println("Running database migrations...")
	if err := repo.Migrate(
		&models.User{},
		&models.Organization{},
		&models.OrganizationMember{},
		&models.Skill{},
		&models.UserSkill{},
		&models.TaskSkill{},
		&models.Project{},
		&models.ProjectMember{},
		&models.Task{},
		&models.TaskDependency{},
		&models.Nudge{},
		&models.NudgeAction{},
		&models.AssignmentSuggestion{},
		&models.WorkloadEntry{},
		&models.Scenario{},
		&models.ScenarioImpactAnalysis{},
	); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Check database health
	if err := repo.Health(); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}
	log.Println("Database connection established")

	// Create repository provider
	repos := repositories.NewProvider(repo.DB())
	log.Println("Repository provider initialized")

	// Setup routes with real services
	router := routes.SetupRoutesWithRepos(repos)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router.GetEngine(),
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
