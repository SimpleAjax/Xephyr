package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/services"
)

// DependencyController handles dependency-related HTTP requests
type DependencyController struct {
	service services.DependencyService
}

// NewDependencyController creates a new dependency controller
func NewDependencyController(service services.DependencyService) *DependencyController {
	return &DependencyController{service: service}
}

// GetTaskDependencies godoc
// @Summary Get task dependencies
// @Description Get all dependencies for a specific task
// @Tags dependencies
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Param includeIndirect query bool false "Include indirect dependencies"
// @Success 200 {object} dto.ApiResponse{data=dto.TaskDependenciesResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /dependencies/tasks/{taskId} [get]
func (c *DependencyController) GetTaskDependencies(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	orgID := ctx.GetString("organizationId")

	var params dto.TaskDependenciesQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	deps, err := c.service.GetTaskDependencies(ctx.Request.Context(), taskID, params.IncludeIndirect, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Task not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(deps, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// CreateDependency godoc
// @Summary Create dependency
// @Description Create a new dependency between tasks
// @Tags dependencies
// @Accept json
// @Produce json
// @Param request body dto.CreateDependencyRequest true "Dependency creation request"
// @Success 201 {object} dto.ApiResponse{data=dto.CreateDependencyResponse}
// @Failure 400 {object} dto.ApiResponse
// @Failure 409 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /dependencies [post]
func (c *DependencyController) CreateDependency(ctx *gin.Context) {
	var req dto.CreateDependencyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	dep, err := c.service.CreateDependency(ctx.Request.Context(), req, orgID)
	if err != nil {
		// Check for circular dependency error
		if circularErr, ok := err.(*services.CircularDependencyError); ok {
			details := map[string]interface{}{
				"cycle": circularErr.Cycle,
			}
			ctx.JSON(http.StatusConflict, dto.NewErrorResponse("CIRCULAR_DEPENDENCY", "This dependency would create a circular reference", details, ctx.GetString("requestId")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewSuccessResponse(dep, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// DeleteDependency godoc
// @Summary Delete dependency
// @Description Remove a dependency
// @Tags dependencies
// @Accept json
// @Produce json
// @Param dependencyId path string true "Dependency ID"
// @Success 204 "No Content"
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /dependencies/{dependencyId} [delete]
func (c *DependencyController) DeleteDependency(ctx *gin.Context) {
	dependencyID := ctx.Param("dependencyId")
	orgID := ctx.GetString("organizationId")

	if err := c.service.DeleteDependency(ctx.Request.Context(), dependencyID, orgID); err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Dependency not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetCriticalPath godoc
// @Summary Get critical path
// @Description Get the critical path for a project
// @Tags dependencies
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 200 {object} dto.ApiResponse{data=dto.CriticalPathResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /dependencies/critical-path/{projectId} [get]
func (c *DependencyController) GetCriticalPath(ctx *gin.Context) {
	projectID := ctx.Param("projectId")
	orgID := ctx.GetString("organizationId")

	path, err := c.service.GetCriticalPath(ctx.Request.Context(), projectID, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Project not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(path, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// ValidateDependency godoc
// @Summary Validate dependency
// @Description Validate a potential dependency without creating it
// @Tags dependencies
// @Accept json
// @Produce json
// @Param request body dto.ValidateDependencyRequest true "Validation request"
// @Success 200 {object} dto.ApiResponse{data=dto.ValidateDependencyResponse}
// @Failure 400 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /dependencies/validate [post]
func (c *DependencyController) ValidateDependency(ctx *gin.Context) {
	var req dto.ValidateDependencyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	result, err := c.service.ValidateDependency(ctx.Request.Context(), req, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(result, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetDependencyGraph godoc
// @Summary Get dependency graph
// @Description Get the dependency graph for visualization
// @Tags dependencies
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 200 {object} dto.ApiResponse{data=dto.DependencyGraphResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /dependencies/graph/{projectId} [get]
func (c *DependencyController) GetDependencyGraph(ctx *gin.Context) {
	projectID := ctx.Param("projectId")
	orgID := ctx.GetString("organizationId")

	graph, err := c.service.GetDependencyGraph(ctx.Request.Context(), projectID, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Project not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(graph, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}
