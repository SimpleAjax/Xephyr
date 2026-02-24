package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/SimpleAjax/Xephyr/internal/controllers"
	"github.com/SimpleAjax/Xephyr/internal/middleware"
	"github.com/SimpleAjax/Xephyr/internal/repositories"
	"github.com/SimpleAjax/Xephyr/internal/services"
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
	projectCtrl *controllers.ProjectController,
	taskCtrl *controllers.TaskController,
	userCtrl *controllers.UserController,
	skillCtrl *controllers.SkillController,
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

	// Swagger documentation (no auth required)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := router.Group("/api/v1")
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
		registerProjectRoutes(v1, projectCtrl)
		registerTaskRoutes(v1, taskCtrl)
		registerUserRoutes(v1, userCtrl)
		registerSkillRoutes(v1, skillCtrl)
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

// registerProjectRoutes registers project module routes
func registerProjectRoutes(rg *gin.RouterGroup, ctrl *controllers.ProjectController) {
	if ctrl == nil {
		return
	}
	projects := rg.Group("/projects")
	{
		// CRUD
		projects.GET("", ctrl.ListProjects)
		projects.POST("", ctrl.CreateProject)
		projects.GET("/:projectId", ctrl.GetProject)
		projects.PATCH("/:projectId", ctrl.UpdateProject)
		projects.DELETE("/:projectId", ctrl.DeleteProject)

		// Team
		projects.GET("/:projectId/team", ctrl.GetProjectTeam)
	}
}

// registerTaskRoutes registers task module routes
func registerTaskRoutes(rg *gin.RouterGroup, ctrl *controllers.TaskController) {
	if ctrl == nil {
		return
	}
	tasks := rg.Group("/tasks")
	{
		// CRUD
		tasks.GET("", ctrl.ListTasks)
		tasks.POST("", ctrl.CreateTask)
		tasks.GET("/:taskId", ctrl.GetTask)
		tasks.PATCH("/:taskId", ctrl.UpdateTask)
		tasks.DELETE("/:taskId", ctrl.DeleteTask)

		// Status updates
		tasks.POST("/:taskId/status", ctrl.UpdateTaskStatus)

		// Assignment
		tasks.POST("/:taskId/assign", ctrl.AssignTask)
	}
}

// registerUserRoutes registers user module routes
func registerUserRoutes(rg *gin.RouterGroup, ctrl *controllers.UserController) {
	if ctrl == nil {
		return
	}
	users := rg.Group("/users")
	{
		// CRUD
		users.GET("", ctrl.ListUsers)
		users.GET("/:userId", ctrl.GetUser)

		// Skills
		users.GET("/:userId/skills", ctrl.GetUserSkills)

		// Workload
		users.GET("/:userId/workload", ctrl.GetUserWorkload)
	}
}

// registerSkillRoutes registers skill module routes
func registerSkillRoutes(rg *gin.RouterGroup, ctrl *controllers.SkillController) {
	if ctrl == nil {
		return
	}
	skills := rg.Group("/skills")
	{
		skills.GET("", ctrl.ListSkills)
		skills.GET("/gaps", ctrl.GetSkillGaps)
	}
}

// SetupRoutesWithRepos creates all services, controllers and routes using real repositories
func SetupRoutesWithRepos(repos *repositories.Provider) *Router {
	// Create real services using repositories
	priorityService := services.NewDummyPriorityService() // Keep dummy for now
	healthService := services.NewRealHealthService(repos)
	nudgeService := services.NewRealNudgeService(repos)
	progressService := services.NewDummyProgressService() // Keep dummy for now
	dependencyService := services.NewDummyDependencyService() // Keep dummy for now
	assignmentService := services.NewDummyAssignmentService() // Keep dummy for now
	scenarioService := services.NewRealScenarioService(repos)
	workloadService := services.NewRealWorkloadService(repos)

	// Create controllers
	priorityCtrl := controllers.NewPriorityController(priorityService)
	healthCtrl := controllers.NewHealthController(healthService)
	nudgeCtrl := controllers.NewNudgeController(nudgeService)
	progressCtrl := controllers.NewProgressController(progressService)
	dependencyCtrl := controllers.NewDependencyController(dependencyService)
	assignmentCtrl := controllers.NewAssignmentController(assignmentService)
	scenarioCtrl := controllers.NewScenarioController(scenarioService)
	workloadCtrl := controllers.NewWorkloadController(workloadService)
	projectCtrl := controllers.NewProjectController(repos)
	taskCtrl := controllers.NewTaskController(repos)
	userCtrl := controllers.NewUserController(repos)
	skillCtrl := controllers.NewSkillController(repos)

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
		projectCtrl,
		taskCtrl,
		userCtrl,
		skillCtrl,
		authMiddleware,
		orgMiddleware,
	)
}

// SetupRoutes is deprecated - use SetupRoutesWithRepos instead
// This function is kept for backward compatibility but will panic
func SetupRoutes() *Router {
	panic("SetupRoutes is deprecated - use SetupRoutesWithRepos with a repository provider instead")
}
