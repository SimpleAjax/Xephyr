package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
	"github.com/xephyr-ai/xephyr-backend/internal/services"
)

// AssignmentController handles assignment-related HTTP requests
type AssignmentController struct {
	service services.AssignmentService
}

// NewAssignmentController creates a new assignment controller
func NewAssignmentController(service services.AssignmentService) *AssignmentController {
	return &AssignmentController{service: service}
}

// GetAssignmentSuggestions godoc
// @Summary Get assignment suggestions
// @Description Get AI-powered assignment suggestions for a task
// @Tags assignments
// @Accept json
// @Produce json
// @Param taskId query string true "Task ID"
// @Param limit query int false "Number of suggestions" default(3)
// @Success 200 {object} dto.ApiResponse{data=dto.AssignmentSuggestionsResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /assignments/suggestions [get]
func (c *AssignmentController) GetAssignmentSuggestions(ctx *gin.Context) {
	var params dto.AssignmentSuggestionsQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	suggestions, err := c.service.GetAssignmentSuggestions(ctx.Request.Context(), params.TaskID, params.Limit, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(suggestions, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// AssignTask godoc
// @Summary Assign task
// @Description Assign a task to a specific person
// @Tags assignments
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Param request body dto.AssignTaskRequest true "Assignment request"
// @Success 200 {object} dto.ApiResponse{data=dto.AssignTaskResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /assignments/tasks/{taskId}/assign [post]
func (c *AssignmentController) AssignTask(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	orgID := ctx.GetString("organizationId")
	assignedByStr := ctx.GetString("userId")
	assignedBy, _ := uuid.Parse(assignedByStr)

	var req dto.AssignTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	result, err := c.service.AssignTask(ctx.Request.Context(), taskID, req, orgID, assignedBy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(result, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// AutoAssignTask godoc
// @Summary Auto-assign task
// @Description Automatically assign a task based on best match
// @Tags assignments
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Param request body dto.AutoAssignTaskRequest true "Auto-assignment request"
// @Success 200 {object} dto.ApiResponse{data=dto.AssignTaskResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /assignments/tasks/{taskId}/auto-assign [post]
func (c *AssignmentController) AutoAssignTask(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	orgID := ctx.GetString("organizationId")
	assignedByStr := ctx.GetString("userId")
	assignedBy, _ := uuid.Parse(assignedByStr)

	var req dto.AutoAssignTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	result, err := c.service.AutoAssignTask(ctx.Request.Context(), taskID, req, orgID, assignedBy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(result, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// CheckCompatibility godoc
// @Summary Check compatibility
// @Description Check person-task compatibility
// @Tags assignments
// @Accept json
// @Produce json
// @Param taskId query string true "Task ID"
// @Param personId query string true "Person ID"
// @Success 200 {object} dto.ApiResponse{data=dto.AssignmentCompatibilityResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /assignments/compatibility [get]
func (c *AssignmentController) CheckCompatibility(ctx *gin.Context) {
	var params dto.CompatibilityQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	compat, err := c.service.CheckCompatibility(ctx.Request.Context(), params.TaskID, params.PersonID, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(compat, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// BulkReassign godoc
// @Summary Bulk reassign
// @Description Reassign multiple tasks in bulk
// @Tags assignments
// @Accept json
// @Produce json
// @Param request body dto.BulkReassignRequest true "Bulk reassign request"
// @Success 200 {object} dto.ApiResponse{data=dto.BulkReassignResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /assignments/bulk-reassign [post]
func (c *AssignmentController) BulkReassign(ctx *gin.Context) {
	orgID := ctx.GetString("organizationId")
	performedByStr := ctx.GetString("userId")
	performedBy, _ := uuid.Parse(performedByStr)

	var req dto.BulkReassignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	result, err := c.service.BulkReassign(ctx.Request.Context(), req, orgID, performedBy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(result, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}
