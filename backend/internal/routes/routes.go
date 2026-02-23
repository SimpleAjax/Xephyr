package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/xephyr-ai/xephyr-backend/internal/controllers"
	"github.com/xephyr-ai/xephyr-backend/internal/middleware"
	"github.com/xephyr-ai/xephyr-backend/internal/services"
)

// Router holds all application routes
type Router struct {
	engine *gin.Engine
}

// NewRouter creates a new router with all routes configured
func NewRouter(
	priorityCtrl *controllers.PriorityController,
	healthCtrl *controllers.HealthController,
	nudgeCtrl *controllers.NudgeController,
	progressCtrl *controllers.ProgressController,
	dependencyCtrl *controllers.DependencyController,
	assignmentCtrl *controllers.AssignmentController,
	scenarioCtrl *controllers.ScenarioController,
	workloadCtrl *controllers.WorkloadController,
	authMiddleware *middleware.AuthMiddleware,
	orgMiddleware *middleware.OrganizationMiddleware,
) *Router {
	router := gin.New()

	// Global middleware
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS(middleware.DefaultCORSConfig()))
	router.Use(middleware.DefaultLogger())

	// Health check (no auth required)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 group
	v1 := router.Group("/v1")
	{
		// Apply auth middleware to all v1 routes
		v1.Use(authMiddleware.Authenticate())
		v1.Use(orgMiddleware.RequireOrganization())

		// Register all module routes
		registerPriorityRoutes(v1, priorityCtrl)
		registerHealthRoutes(v1, healthCtrl)
		registerNudgeRoutes(v1, nudgeCtrl)
		registerProgressRoutes(v1, progressCtrl)
		registerDependencyRoutes(v1, dependencyCtrl)
		registerAssignmentRoutes(v1, assignmentCtrl)
		registerScenarioRoutes(v1, scenarioCtrl)
		registerWorkloadRoutes(v1, workloadCtrl)
	}

	// Handle 404s
	router.NoRoute(middleware.NotFoundHandler())
	router.NoMethod(middleware.MethodNotAllowedHandler())

	return &Router{engine: router}
}

// GetEngine returns the gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// registerPriorityRoutes registers priority module routes
func registerPriorityRoutes(rg *gin.RouterGroup, ctrl *controllers.PriorityController) {
	priorities := rg.Group("/priorities")
	{
		// Task priorities
		priorities.GET("/tasks/:taskId", ctrl.GetTaskPriority)
		priorities.POST("/tasks/bulk", ctrl.GetBulkTaskPriorities)

		// Project rankings
		priorities.GET("/projects/:projectId/ranking", ctrl.GetProjectTaskRanking)

		// Recalculation
		priorities.POST("/recalculate", ctrl.RecalculatePriorities)
	}
}

// registerHealthRoutes registers health module routes
func registerHealthRoutes(rg *gin.RouterGroup, ctrl *controllers.HealthController) {
	health := rg.Group("/health")
	{
		// Portfolio
		health.GET("/portfolio", ctrl.GetPortfolioHealth)

		// Projects
		health.GET("/projects", ctrl.GetBulkProjectHealth)
		health.GET("/projects/:projectId", ctrl.GetProjectHealth)

		// Trends
		health.GET("/trends", ctrl.GetHealthTrends)
	}
}

// registerNudgeRoutes registers nudge module routes
func registerNudgeRoutes(rg *gin.RouterGroup, ctrl *controllers.NudgeController) {
	nudges := rg.Group("/nudges")
	{
		// CRUD operations
		nudges.GET("", ctrl.ListNudges)
		nudges.GET("/:nudgeId", ctrl.GetNudge)
		nudges.PATCH("/:nudgeId/status", ctrl.UpdateNudgeStatus)

		// Actions
		nudges.POST("/:nudgeId/actions", ctrl.TakeNudgeAction)

		// Generation
		nudges.POST("/generate", ctrl.GenerateNudges)

		// Stats
		nudges.GET("/stats", ctrl.GetNudgeStats)
	}
}

// registerProgressRoutes registers progress module routes
func registerProgressRoutes(rg *gin.RouterGroup, ctrl *controllers.ProgressController) {
	progress := rg.Group("/progress")
	{
		// Projects
		progress.GET("/projects/:projectId", ctrl.GetProjectProgress)

		// Tasks
		progress.GET("/tasks/:taskId", ctrl.GetTaskProgress)
		progress.POST("/tasks/:taskId/update", ctrl.UpdateTaskProgress)

		// Rollups
		progress.GET("/rollups/:projectId", ctrl.GetProjectRollup)
	}
}

// registerDependencyRoutes registers dependency module routes
func registerDependencyRoutes(rg *gin.RouterGroup, ctrl *controllers.DependencyController) {
	dependencies := rg.Group("/dependencies")
	{
		// Task dependencies
		dependencies.GET("/tasks/:taskId", ctrl.GetTaskDependencies)

		// CRUD
		dependencies.POST("", ctrl.CreateDependency)
		dependencies.DELETE("/:dependencyId", ctrl.DeleteDependency)

		// Critical path
		dependencies.GET("/critical-path/:projectId", ctrl.GetCriticalPath)

		// Validation
		dependencies.POST("/validate", ctrl.ValidateDependency)

		// Graph
		dependencies.GET("/graph/:projectId", ctrl.GetDependencyGraph)
	}
}

// registerAssignmentRoutes registers assignment module routes
func registerAssignmentRoutes(rg *gin.RouterGroup, ctrl *controllers.AssignmentController) {
	assignments := rg.Group("/assignments")
	{
		// Suggestions
		assignments.GET("/suggestions", ctrl.GetAssignmentSuggestions)

		// Assignment operations
		assignments.POST("/tasks/:taskId/assign", ctrl.AssignTask)
		assignments.POST("/tasks/:taskId/auto-assign", ctrl.AutoAssignTask)

		// Compatibility
		assignments.GET("/compatibility", ctrl.CheckCompatibility)

		// Bulk operations
		assignments.POST("/bulk-reassign", ctrl.BulkReassign)
	}
}

// registerScenarioRoutes registers scenario module routes
func registerScenarioRoutes(rg *gin.RouterGroup, ctrl *controllers.ScenarioController) {
	scenarios := rg.Group("/scenarios")
	{
		// CRUD
		scenarios.POST("", ctrl.CreateScenario)
		scenarios.GET("", ctrl.ListScenarios)
		scenarios.GET("/:scenarioId", ctrl.GetScenario)

		// Simulation
		scenarios.POST("/:scenarioId/simulate", ctrl.SimulateScenario)

		// Actions
		scenarios.POST("/:scenarioId/apply", ctrl.ApplyScenario)
		scenarios.POST("/:scenarioId/reject", ctrl.RejectScenario)
		scenarios.PATCH("/:scenarioId/modify", ctrl.ModifyScenario)
	}
}

// registerWorkloadRoutes registers workload module routes
func registerWorkloadRoutes(rg *gin.RouterGroup, ctrl *controllers.WorkloadController) {
	workload := rg.Group("/workload")
	{
		// Team workload
		workload.GET("/team", ctrl.GetTeamWorkload)

		// Individual
		workload.GET("/people/:personId", ctrl.GetIndividualWorkload)

		// Forecast
		workload.GET("/forecast", ctrl.GetWorkloadForecast)

		// Analytics
		workload.GET("/analytics", ctrl.GetWorkloadAnalytics)

		// Rebalancing
		workload.POST("/rebalance", ctrl.GetRebalanceSuggestions)
	}
}

// SetupRoutes is a convenience function that creates all services, controllers and routes
func SetupRoutes() *Router {
	// Create services
	priorityService := services.NewDummyPriorityService()
	healthService := services.NewDummyHealthService()
	nudgeService := services.NewDummyNudgeService()
	progressService := services.NewDummyProgressService()
	dependencyService := services.NewDummyDependencyService()
	assignmentService := services.NewDummyAssignmentService()
	scenarioService := services.NewDummyScenarioService()
	workloadService := services.NewDummyWorkloadService()

	// Create controllers
	priorityCtrl := controllers.NewPriorityController(priorityService)
	healthCtrl := controllers.NewHealthController(healthService)
	nudgeCtrl := controllers.NewNudgeController(nudgeService)
	progressCtrl := controllers.NewProgressController(progressService)
	dependencyCtrl := controllers.NewDependencyController(dependencyService)
	assignmentCtrl := controllers.NewAssignmentController(assignmentService)
	scenarioCtrl := controllers.NewScenarioController(scenarioService)
	workloadCtrl := controllers.NewWorkloadController(workloadService)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware("dummy-secret")
	orgMiddleware := middleware.NewOrganizationMiddleware()

	// Create and return router
	return NewRouter(
		priorityCtrl,
		healthCtrl,
		nudgeCtrl,
		progressCtrl,
		dependencyCtrl,
		assignmentCtrl,
		scenarioCtrl,
		workloadCtrl,
		authMiddleware,
		orgMiddleware,
	)
}
