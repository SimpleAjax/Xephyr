package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/services"
)

// ProgressController handles progress-related HTTP requests
type ProgressController struct {
	service services.ProgressService
}

// NewProgressController creates a new progress controller
func NewProgressController(service services.ProgressService) *ProgressController {
	return &ProgressController{service: service}
}

// GetProjectProgress godoc
// @Summary Get project progress
// @Description Get progress information for a project
// @Tags progress
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 200 {object} dto.ApiResponse{data=dto.ProjectProgressResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /progress/projects/{projectId} [get]
func (c *ProgressController) GetProjectProgress(ctx *gin.Context) {
	projectID := ctx.Param("projectId")
	orgID := ctx.GetString("organizationId")

	progress, err := c.service.GetProjectProgress(ctx.Request.Context(), projectID, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Project not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(progress, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetTaskProgress godoc
// @Summary Get task progress
// @Description Get detailed progress for a task
// @Tags progress
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} dto.ApiResponse{data=dto.TaskProgressResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /progress/tasks/{taskId} [get]
func (c *ProgressController) GetTaskProgress(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	orgID := ctx.GetString("organizationId")

	progress, err := c.service.GetTaskProgress(ctx.Request.Context(), taskID, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Task not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(progress, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// UpdateTaskProgress godoc
// @Summary Update task progress
// @Description Update the progress of a task
// @Tags progress
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Param request body dto.UpdateTaskProgressRequest true "Progress update request"
// @Success 200 {object} dto.ApiResponse{data=dto.TaskProgressUpdateResponse}
// @Failure 400 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /progress/tasks/{taskId}/update [post]
func (c *ProgressController) UpdateTaskProgress(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	orgID := ctx.GetString("organizationId")
	userID := ctx.GetString("userId")

	var req dto.UpdateTaskProgressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	result, err := c.service.UpdateTaskProgress(ctx.Request.Context(), taskID, req, orgID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(result, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetProjectRollup godoc
// @Summary Get project rollup
// @Description Get hierarchical progress rollup for a project
// @Tags progress
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 200 {object} dto.ApiResponse{data=dto.ProjectRollupResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /progress/rollups/{projectId} [get]
func (c *ProgressController) GetProjectRollup(ctx *gin.Context) {
	projectID := ctx.Param("projectId")
	orgID := ctx.GetString("organizationId")

	rollup, err := c.service.GetProjectRollup(ctx.Request.Context(), projectID, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Project not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(rollup, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}
