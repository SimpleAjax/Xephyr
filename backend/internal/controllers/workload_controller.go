package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/services"
)

// WorkloadController handles workload-related HTTP requests
type WorkloadController struct {
	service services.WorkloadService
}

// NewWorkloadController creates a new workload controller
func NewWorkloadController(service services.WorkloadService) *WorkloadController {
	return &WorkloadController{service: service}
}

// GetTeamWorkload godoc
// @Summary Get team workload
// @Description Get team workload overview
// @Tags workload
// @Accept json
// @Produce json
// @Param week query string false "Week starting date"
// @Param includeForecast query bool false "Include forecast"
// @Success 200 {object} dto.ApiResponse{data=dto.TeamWorkloadResponse}
// @Security BearerAuth
// @Router /workload/team [get]
func (c *WorkloadController) GetTeamWorkload(ctx *gin.Context) {
	var params dto.TeamWorkloadQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	workload, err := c.service.GetTeamWorkload(ctx.Request.Context(), params.Week, params.IncludeForecast, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(workload, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetIndividualWorkload godoc
// @Summary Get individual workload
// @Description Get workload for a specific person
// @Tags workload
// @Accept json
// @Produce json
// @Param personId path string true "Person ID"
// @Success 200 {object} dto.ApiResponse{data=dto.IndividualWorkloadResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /workload/people/{personId} [get]
func (c *WorkloadController) GetIndividualWorkload(ctx *gin.Context) {
	personID := ctx.Param("personId")
	orgID := ctx.GetString("organizationId")

	workload, err := c.service.GetIndividualWorkload(ctx.Request.Context(), personID, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Person not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(workload, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetWorkloadForecast godoc
// @Summary Get workload forecast
// @Description Get workload forecast for a person
// @Tags workload
// @Accept json
// @Produce json
// @Param personId query string true "Person ID"
// @Param weeks query int false "Number of weeks" default(8)
// @Success 200 {object} dto.ApiResponse{data=dto.WorkloadForecastResponse}
// @Failure 400 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /workload/forecast [get]
func (c *WorkloadController) GetWorkloadForecast(ctx *gin.Context) {
	var params dto.WorkloadForecastQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	forecast, err := c.service.GetWorkloadForecast(ctx.Request.Context(), params.PersonID, params.Weeks, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(forecast, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetWorkloadAnalytics godoc
// @Summary Get workload analytics
// @Description Get workload analytics for the organization
// @Tags workload
// @Accept json
// @Produce json
// @Param period query string false "Time period" default(30d)
// @Success 200 {object} dto.ApiResponse{data=dto.WorkloadAnalyticsResponse}
// @Security BearerAuth
// @Router /workload/analytics [get]
func (c *WorkloadController) GetWorkloadAnalytics(ctx *gin.Context) {
	var params dto.WorkloadAnalyticsQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	analytics, err := c.service.GetWorkloadAnalytics(ctx.Request.Context(), params.Period, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(analytics, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetRebalanceSuggestions godoc
// @Summary Get rebalance suggestions
// @Description Get suggestions for rebalancing workload
// @Tags workload
// @Accept json
// @Produce json
// @Param request body dto.RebalanceWorkloadRequest true "Rebalance request"
// @Success 200 {object} dto.ApiResponse{data=dto.RebalanceWorkloadResponse}
// @Failure 400 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /workload/rebalance [post]
func (c *WorkloadController) GetRebalanceSuggestions(ctx *gin.Context) {
	var req dto.RebalanceWorkloadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	suggestions, err := c.service.GetRebalanceSuggestions(ctx.Request.Context(), req, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(suggestions, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}
