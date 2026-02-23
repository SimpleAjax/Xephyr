package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
	"github.com/xephyr-ai/xephyr-backend/internal/services"
)

// PriorityController handles priority-related HTTP requests
type PriorityController struct {
	service services.PriorityService
}

// NewPriorityController creates a new priority controller
func NewPriorityController(service services.PriorityService) *PriorityController {
	return &PriorityController{service: service}
}

// GetTaskPriority godoc
// @Summary Get single task priority
// @Description Get the priority score and breakdown for a specific task
// @Tags priorities
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} dto.ApiResponse{data=dto.TaskPriorityResponse}
// @Failure 404 {object} dto.ApiResponse
// @Router /priorities/tasks/{taskId} [get]
func (c *PriorityController) GetTaskPriority(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	orgID := ctx.GetString("organizationId")

	priority, err := c.service.GetTaskPriority(ctx.Request.Context(), taskID, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Task not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(priority, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetBulkTaskPriorities godoc
// @Summary Get priorities for multiple tasks
// @Description Get priority scores for multiple tasks in a single request
// @Tags priorities
// @Accept json
// @Produce json
// @Param request body dto.BulkPriorityRequest true "Bulk priority request"
// @Success 200 {object} dto.ApiResponse{data=dto.BulkPriorityResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /priorities/tasks/bulk [post]
func (c *PriorityController) GetBulkTaskPriorities(ctx *gin.Context) {
	var req dto.BulkPriorityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	priorities, err := c.service.GetBulkTaskPriorities(ctx.Request.Context(), req, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(priorities, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetProjectTaskRanking godoc
// @Summary Get task ranking for project
// @Description Get prioritized task ranking for a specific project
// @Tags priorities
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Param status query string false "Filter by status"
// @Param assigneeId query string false "Filter by assignee"
// @Param minScore query int false "Minimum priority score"
// @Param limit query int false "Max results (default 50)" default(50)
// @Param offset query int false "Pagination offset" default(0)
// @Success 200 {object} dto.ApiResponse{data=dto.ProjectTaskRankingResponse,meta=dto.ResponseMeta}
// @Failure 404 {object} dto.ApiResponse
// @Router /priorities/projects/{projectId}/ranking [get]
func (c *PriorityController) GetProjectTaskRanking(ctx *gin.Context) {
	projectID := ctx.Param("projectId")
	orgID := ctx.GetString("organizationId")

	var params dto.ProjectRankingQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ranking, err := c.service.GetProjectTaskRanking(ctx.Request.Context(), projectID, params, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Project not found", nil, ctx.GetString("requestId")))
		return
	}

	meta := dto.ResponseMeta{
		Page:      &params.Offset,
		PerPage:   &params.Limit,
		Total:     &ranking.Total,
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(ranking, meta))
}

// RecalculatePriorities godoc
// @Summary Trigger priority recalculation
// @Description Trigger recalculation of priorities for a scope
// @Tags priorities
// @Accept json
// @Produce json
// @Param request body dto.RecalculatePriorityRequest true "Recalculation request"
// @Success 200 {object} dto.ApiResponse{data=dto.RecalculatePrioritySyncResponse}
// @Success 202 {object} dto.ApiResponse{data=dto.RecalculatePriorityAsyncResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /priorities/recalculate [post]
func (c *PriorityController) RecalculatePriorities(ctx *gin.Context) {
	var req dto.RecalculatePriorityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	syncResp, asyncResp, err := c.service.RecalculatePriorities(ctx.Request.Context(), req, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	if req.Async {
		ctx.JSON(http.StatusAccepted, dto.NewSuccessResponse(asyncResp, dto.ResponseMeta{
			Timestamp: getTimestamp(),
			RequestID: ctx.GetString("requestId"),
		}))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(syncResp, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// Helper function
func getTimestamp() string {
	return uuid.New().String()[:8]
}
