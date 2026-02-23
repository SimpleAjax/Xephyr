package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
	"github.com/xephyr-ai/xephyr-backend/internal/services"
)

// ScenarioController handles scenario-related HTTP requests
type ScenarioController struct {
	service services.ScenarioService
}

// NewScenarioController creates a new scenario controller
func NewScenarioController(service services.ScenarioService) *ScenarioController {
	return &ScenarioController{service: service}
}

// CreateScenario godoc
// @Summary Create scenario
// @Description Create a new what-if scenario
// @Tags scenarios
// @Accept json
// @Produce json
// @Param request body dto.CreateScenarioRequest true "Scenario creation request"
// @Success 201 {object} dto.ApiResponse{data=dto.ScenarioResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /scenarios [post]
func (c *ScenarioController) CreateScenario(ctx *gin.Context) {
	var req dto.CreateScenarioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	createdByStr := ctx.GetString("userId")
	createdBy, _ := uuid.Parse(createdByStr)

	scenario, err := c.service.CreateScenario(ctx.Request.Context(), req, orgID, createdBy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewSuccessResponse(scenario, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// ListScenarios godoc
// @Summary List scenarios
// @Description Get a list of scenarios
// @Tags scenarios
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Param limit query int false "Limit results" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} dto.ApiResponse{data=dto.ScenarioListResponse,meta=dto.ResponseMeta}
// @Router /scenarios [get]
func (c *ScenarioController) ListScenarios(ctx *gin.Context) {
	var params dto.ScenarioListQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	orgID := ctx.GetString("organizationId")
	scenarios, err := c.service.ListScenarios(ctx.Request.Context(), params, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	meta := dto.ResponseMeta{
		Page:      &params.Offset,
		PerPage:   &params.Limit,
		Total:     &scenarios.Total,
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(scenarios, meta))
}

// GetScenario godoc
// @Summary Get scenario
// @Description Get detailed information about a scenario
// @Tags scenarios
// @Accept json
// @Produce json
// @Param scenarioId path string true "Scenario ID"
// @Success 200 {object} dto.ApiResponse{data=dto.ScenarioDetailResponse}
// @Failure 404 {object} dto.ApiResponse
// @Router /scenarios/{scenarioId} [get]
func (c *ScenarioController) GetScenario(ctx *gin.Context) {
	scenarioID := ctx.Param("scenarioId")
	orgID := ctx.GetString("organizationId")

	scenario, err := c.service.GetScenario(ctx.Request.Context(), scenarioID, orgID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Scenario not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(scenario, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// SimulateScenario godoc
// @Summary Simulate scenario
// @Description Run a simulation for a scenario
// @Tags scenarios
// @Accept json
// @Produce json
// @Param scenarioId path string true "Scenario ID"
// @Param request body dto.SimulateScenarioRequest true "Simulation request"
// @Success 200 {object} dto.ApiResponse{data=dto.SimulateScenarioResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /scenarios/{scenarioId}/simulate [post]
func (c *ScenarioController) SimulateScenario(ctx *gin.Context) {
	scenarioID := ctx.Param("scenarioId")
	orgID := ctx.GetString("organizationId")

	var req dto.SimulateScenarioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	result, err := c.service.SimulateScenario(ctx.Request.Context(), scenarioID, req, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(result, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// ApplyScenario godoc
// @Summary Apply scenario
// @Description Apply a scenario's changes
// @Tags scenarios
// @Accept json
// @Produce json
// @Param scenarioId path string true "Scenario ID"
// @Param request body dto.ApplyScenarioRequest true "Apply request"
// @Success 200 {object} dto.ApiResponse{data=dto.ApplyScenarioResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /scenarios/{scenarioId}/apply [post]
func (c *ScenarioController) ApplyScenario(ctx *gin.Context) {
	scenarioID := ctx.Param("scenarioId")
	orgID := ctx.GetString("organizationId")
	appliedByStr := ctx.GetString("userId")
	appliedBy, _ := uuid.Parse(appliedByStr)

	var req dto.ApplyScenarioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	result, err := c.service.ApplyScenario(ctx.Request.Context(), scenarioID, req, orgID, appliedBy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(result, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// RejectScenario godoc
// @Summary Reject scenario
// @Description Reject a scenario
// @Tags scenarios
// @Accept json
// @Produce json
// @Param scenarioId path string true "Scenario ID"
// @Success 200 {object} dto.ApiResponse{data=dto.RejectScenarioResponse}
// @Failure 404 {object} dto.ApiResponse
// @Router /scenarios/{scenarioId}/reject [post]
func (c *ScenarioController) RejectScenario(ctx *gin.Context) {
	scenarioID := ctx.Param("scenarioId")
	orgID := ctx.GetString("organizationId")
	rejectedByStr := ctx.GetString("userId")
	rejectedBy, _ := uuid.Parse(rejectedByStr)

	result, err := c.service.RejectScenario(ctx.Request.Context(), scenarioID, orgID, rejectedBy)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Scenario not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(result, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// ModifyScenario godoc
// @Summary Modify scenario
// @Description Modify an existing scenario
// @Tags scenarios
// @Accept json
// @Produce json
// @Param scenarioId path string true "Scenario ID"
// @Param request body dto.ModifyScenarioRequest true "Modification request"
// @Success 200 {object} dto.ApiResponse{data=dto.ScenarioResponse}
// @Failure 400 {object} dto.ApiResponse
// @Router /scenarios/{scenarioId}/modify [patch]
func (c *ScenarioController) ModifyScenario(ctx *gin.Context) {
	scenarioID := ctx.Param("scenarioId")
	orgID := ctx.GetString("organizationId")

	var req dto.ModifyScenarioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	scenario, err := c.service.ModifyScenario(ctx.Request.Context(), scenarioID, req, orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(scenario, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}
