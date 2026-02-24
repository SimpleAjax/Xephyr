package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/services"
)

// NudgeController handles nudge-related HTTP requests
type NudgeController struct {
	service services.NudgeService
}

// NewNudgeController creates a new nudge controller
func NewNudgeController(service services.NudgeService) *NudgeController {
	return &NudgeController{service: service}
}

// ListNudges godoc
// @Summary List nudges
// @Description Get a list of nudges with optional filters
// @Tags nudges
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Param severity query string false "Filter by severity"
// @Param type query string false "Filter by type"
// @Param projectId query string false "Filter by project"
// @Param personId query string false "Filter by person"
// @Param limit query int false "Limit results" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} dto.ApiResponse{data=dto.NudgeListResponse,meta=dto.ResponseMeta}
// @Security BearerAuth
// @Router /nudges [get]
func (c *NudgeController) ListNudges(ctx *gin.Context) {
	var params dto.NudgeListQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	userIDStr := ctx.GetString("userId")
	userID, _ := uuid.Parse(userIDStr)

	nudges, err := c.service.ListNudges(ctx.Request.Context(), params, orgID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	meta := dto.ResponseMeta{
		Page:      &params.Offset,
		PerPage:   &params.Limit,
		Total:     &nudges.Summary.Total,
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nudges, meta))
}

// GetNudge godoc
// @Summary Get single nudge
// @Description Get detailed information about a specific nudge
// @Tags nudges
// @Accept json
// @Produce json
// @Param nudgeId path string true "Nudge ID"
// @Success 200 {object} dto.ApiResponse{data=dto.NudgeDetailResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /nudges/{nudgeId} [get]
func (c *NudgeController) GetNudge(ctx *gin.Context) {
	nudgeID := ctx.Param("nudgeId")
	orgID := ctx.GetString("organizationId")

	nudge, err := c.service.GetNudge(ctx.Request.Context(), nudgeID, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Nudge not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nudge, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// TakeNudgeAction godoc
// @Summary Take action on nudge
// @Description Perform an action on a nudge (accept, dismiss, etc.)
// @Tags nudges
// @Accept json
// @Produce json
// @Param nudgeId path string true "Nudge ID"
// @Param request body dto.NudgeActionRequest true "Action request"
// @Success 200 {object} dto.ApiResponse{data=dto.NudgeActionResponse}
// @Failure 400 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /nudges/{nudgeId}/actions [post]
func (c *NudgeController) TakeNudgeAction(ctx *gin.Context) {
	nudgeID := ctx.Param("nudgeId")
	orgID := ctx.GetString("organizationId")
	userIDStr := ctx.GetString("userId")
	userID, _ := uuid.Parse(userIDStr)

	var req dto.NudgeActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	result, err := c.service.TakeNudgeAction(ctx.Request.Context(), nudgeID, req, orgID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(result, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// UpdateNudgeStatus godoc
// @Summary Update nudge status
// @Description Update the status of a nudge
// @Tags nudges
// @Accept json
// @Produce json
// @Param nudgeId path string true "Nudge ID"
// @Param request body dto.UpdateNudgeStatusRequest true "Status update request"
// @Success 200 {object} dto.ApiResponse{data=dto.NudgeResponse}
// @Failure 400 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /nudges/{nudgeId}/status [patch]
func (c *NudgeController) UpdateNudgeStatus(ctx *gin.Context) {
	nudgeID := ctx.Param("nudgeId")
	orgID := ctx.GetString("organizationId")

	var req dto.UpdateNudgeStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	nudge, err := c.service.UpdateNudgeStatus(ctx.Request.Context(), nudgeID, req.Status, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nudge, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GenerateNudges godoc
// @Summary Generate nudges
// @Description Trigger manual nudge generation
// @Tags nudges
// @Accept json
// @Produce json
// @Param request body dto.GenerateNudgesRequest true "Generation request"
// @Success 202 {object} dto.ApiResponse{data=dto.IDResponse}
// @Failure 400 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /nudges/generate [post]
func (c *NudgeController) GenerateNudges(ctx *gin.Context) {
	var req dto.GenerateNudgesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	jobID, err := c.service.GenerateNudges(ctx.Request.Context(), req, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusAccepted, dto.NewSuccessResponse(dto.IDResponse{ID: jobID}, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetNudgeStats godoc
// @Summary Get nudge statistics
// @Description Get statistics about nudges
// @Tags nudges
// @Accept json
// @Produce json
// @Param period query string false "Time period" default(30d)
// @Success 200 {object} dto.ApiResponse{data=dto.NudgeStatsResponse}
// @Security BearerAuth
// @Router /nudges/stats [get]
func (c *NudgeController) GetNudgeStats(ctx *gin.Context) {
	var params dto.NudgeStatsQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	stats, err := c.service.GetNudgeStats(ctx.Request.Context(), params.Period, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(stats, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}
