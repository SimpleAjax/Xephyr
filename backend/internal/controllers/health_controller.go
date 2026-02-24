package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/services"
)

// HealthController handles health-related HTTP requests
type HealthController struct {
	service services.HealthService
}

// NewHealthController creates a new health controller
func NewHealthController(service services.HealthService) *HealthController {
	return &HealthController{service: service}
}

// GetPortfolioHealth godoc
// @Summary Get portfolio health overview
// @Description Get overall portfolio health score and project summaries
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=dto.PortfolioHealthResponse}
// @Security BearerAuth
// @Router /health/portfolio [get]
func (c *HealthController) GetPortfolioHealth(ctx *gin.Context) {
	orgID := ctx.GetString("organizationId")

	health, err := c.service.GetPortfolioHealth(ctx.Request.Context(), orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(health, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetProjectHealth godoc
// @Summary Get project health details
// @Description Get detailed health metrics for a specific project
// @Tags health
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Param includeBreakdown query bool false "Include health breakdown"
// @Success 200 {object} dto.ApiResponse{data=dto.ProjectHealthResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /health/projects/{projectId} [get]
func (c *HealthController) GetProjectHealth(ctx *gin.Context) {
	projectID := ctx.Param("projectId")
	orgID := ctx.GetString("organizationId")

	var params dto.ProjectHealthQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	health, err := c.service.GetProjectHealth(ctx.Request.Context(), projectID, params.IncludeBreakdown, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Project not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(health, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetBulkProjectHealth godoc
// @Summary Get bulk project health
// @Description Get health summaries for multiple projects
// @Tags health
// @Accept json
// @Produce json
// @Param projectIds query []string false "Project IDs"
// @Success 200 {object} dto.ApiResponse{data=[]dto.ProjectHealthSummary}
// @Security BearerAuth
// @Router /health/projects [get]
func (c *HealthController) GetBulkProjectHealth(ctx *gin.Context) {
	projectIDs := ctx.QueryArray("projectIds")
	orgID := ctx.GetString("organizationId")

	// If no project IDs provided, get all projects from portfolio health
	if len(projectIDs) == 0 {
		portfolioHealth, err := c.service.GetPortfolioHealth(ctx.Request.Context(), orgID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
			return
		}
		ctx.JSON(http.StatusOK, dto.NewSuccessResponse(portfolioHealth.Projects, dto.ResponseMeta{
			Timestamp: getTimestamp(),
			RequestID: ctx.GetString("requestId"),
		}))
		return
	}

	healths, err := c.service.GetBulkProjectHealth(ctx.Request.Context(), projectIDs, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(healths, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetHealthTrends godoc
// @Summary Get health trends
// @Description Get health trend data over time
// @Tags health
// @Accept json
// @Produce json
// @Param projectId query string true "Project ID"
// @Param days query int false "Number of days" default(30)
// @Success 200 {object} dto.ApiResponse{data=dto.HealthTrendsResponse}
// @Failure 400 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /health/trends [get]
func (c *HealthController) GetHealthTrends(ctx *gin.Context) {
	var params dto.HealthTrendsQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	trends, err := c.service.GetHealthTrends(ctx.Request.Context(), params.ProjectID, params.Days, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(trends, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}
