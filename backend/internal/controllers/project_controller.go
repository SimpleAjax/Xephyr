package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/internal/repositories"
)

// ProjectController handles project CRUD HTTP requests
type ProjectController struct {
	repos repositories.Repositories
}

// NewProjectController creates a new project controller
func NewProjectController(repos repositories.Repositories) *ProjectController {
	return &ProjectController{repos: repos}
}

// ListProjects godoc
// @Summary List all projects
// @Description Get a list of projects in the organization
// @Tags projects
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Param limit query int false "Limit results" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} dto.ApiResponse{data=ProjectListResponse}
// @Security BearerAuth
// @Router /projects [get]
func (c *ProjectController) ListProjects(ctx *gin.Context) {
	orgID := ctx.GetString("organizationId")
	log.Printf("[ProjectController] ListProjects called with orgID from context: %s", orgID)
	log.Printf("[ProjectController] X-Organization-Id header: %s", ctx.GetHeader("X-Organization-Id"))
	
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		log.Printf("[ProjectController] Error parsing orgID: %v", err)
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid organization ID", nil, ctx.GetString("requestId")))
		return
	}

	status := ctx.Query("status")

	var params repositories.ListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		params = repositories.DefaultListParams()
	}
	
	// Ensure valid pagination params
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}
	
	log.Printf("[ProjectController] Using params: limit=%d, offset=%d, sort=%s %s", params.Limit, params.Offset, params.SortBy, params.SortOrder)

	var projects []models.Project
	var total int64

	if status != "" {
		projectStatus := models.ProjectStatus(status)
		projects, total, err = c.repos.GetProject().ListByStatus(ctx.Request.Context(), projectStatus, params)
	} else {
		projects, total, err = c.repos.GetProject().ListByOrganization(ctx.Request.Context(), orgUUID, params)
	}

	if err != nil {
		log.Printf("[ProjectController] Error from repository: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	log.Printf("[ProjectController] Got %d projects from repository, total: %d", len(projects), total)

	response := ProjectListResponse{
		Projects: make([]ProjectResponse, 0, len(projects)),
		Total:    int(total),
	}

	for _, p := range projects {
		response.Projects = append(response.Projects, toProjectResponse(&p))
	}

	log.Printf("[ProjectController] Returning %d projects in response", len(response.Projects))

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(response, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetProject godoc
// @Summary Get a project by ID
// @Description Get detailed information about a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 200 {object} dto.ApiResponse{data=ProjectResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /projects/{projectId} [get]
func (c *ProjectController) GetProject(ctx *gin.Context) {
	projectID := ctx.Param("projectId")
	projUUID, err := uuid.Parse(projectID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid project ID", nil, ctx.GetString("requestId")))
		return
	}

	project, err := c.repos.GetProject().GetByID(ctx.Request.Context(), projUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Project not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(toProjectResponse(project), dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project in the organization
// @Tags projects
// @Accept json
// @Produce json
// @Param request body CreateProjectRequest true "Project creation request"
// @Success 201 {object} dto.ApiResponse{data=ProjectResponse}
// @Failure 400 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /projects [post]
func (c *ProjectController) CreateProject(ctx *gin.Context) {
	orgID := ctx.GetString("organizationId")
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid organization ID", nil, ctx.GetString("requestId")))
		return
	}

	var req CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	project := &models.Project{
		OrganizationID: orgUUID,
		Name:           req.Name,
		Description:    req.Description,
		Status:         models.ProjectActive,
		Priority:       req.Priority,
		HealthScore:    100,
		Progress:       0,
	}

	if req.StartDate != nil {
		project.StartDate = req.StartDate
	}
	if req.TargetEndDate != nil {
		project.TargetEndDate = req.TargetEndDate
	}

	if err := c.repos.GetProject().Create(ctx.Request.Context(), project); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewSuccessResponse(toProjectResponse(project), dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update an existing project
// @Tags projects
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Param request body UpdateProjectRequest true "Project update request"
// @Success 200 {object} dto.ApiResponse{data=ProjectResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /projects/{projectId} [patch]
func (c *ProjectController) UpdateProject(ctx *gin.Context) {
	projectID := ctx.Param("projectId")
	projUUID, err := uuid.Parse(projectID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid project ID", nil, ctx.GetString("requestId")))
		return
	}

	project, err := c.repos.GetProject().GetByID(ctx.Request.Context(), projUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Project not found", nil, ctx.GetString("requestId")))
		return
	}

	var req UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	// Update fields
	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.Status != "" {
		project.Status = models.ProjectStatus(req.Status)
	}
	if req.Priority != 0 {
		project.Priority = req.Priority
	}
	if req.HealthScore != nil {
		project.HealthScore = *req.HealthScore
	}
	if req.Progress != nil {
		project.Progress = *req.Progress
	}
	if req.TargetEndDate != nil {
		project.TargetEndDate = req.TargetEndDate
	}

	if err := c.repos.GetProject().Update(ctx.Request.Context(), project); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(toProjectResponse(project), dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Soft-delete a project
// @Tags projects
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 204 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /projects/{projectId} [delete]
func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	projectID := ctx.Param("projectId")
	projUUID, err := uuid.Parse(projectID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid project ID", nil, ctx.GetString("requestId")))
		return
	}

	if err := c.repos.GetProject().Delete(ctx.Request.Context(), projUUID); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetProjectTeam godoc
// @Summary Get project team members
// @Description Get all team members assigned to a project
// @Tags projects
// @Accept json
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 200 {object} dto.ApiResponse{data=ProjectTeamResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /projects/{projectId}/team [get]
func (c *ProjectController) GetProjectTeam(ctx *gin.Context) {
	projectID := ctx.Param("projectId")
	projUUID, err := uuid.Parse(projectID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid project ID", nil, ctx.GetString("requestId")))
		return
	}

	project, err := c.repos.GetProject().GetByID(ctx.Request.Context(), projUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Project not found", nil, ctx.GetString("requestId")))
		return
	}

	members := make([]TeamMemberResponse, 0, len(project.Members))
	for _, m := range project.Members {
		members = append(members, TeamMemberResponse{
			ID:   m.UserID.String(),
			Name: m.User.Name,
			Role: m.Role,
		})
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(ProjectTeamResponse{
		Users: members,
	}, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// Request/Response types

type CreateProjectRequest struct {
	Name          string     `json:"name" binding:"required"`
	Description   string     `json:"description"`
	Priority      int        `json:"priority"`
	StartDate     *time.Time `json:"startDate"`
	TargetEndDate *time.Time `json:"targetEndDate"`
	Budget        float64    `json:"budget"`
}

type UpdateProjectRequest struct {
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Status        string     `json:"status"`
	Priority      int        `json:"priority"`
	HealthScore   *int       `json:"healthScore"`
	Progress      *int       `json:"progress"`
	TargetEndDate *time.Time `json:"targetEndDate"`
}

type ProjectResponse struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Status        string     `json:"status"`
	Priority      int        `json:"priority"`
	HealthScore   int        `json:"healthScore"`
	Progress      int        `json:"progress"`
	StartDate     *time.Time `json:"startDate,omitempty"`
	TargetEndDate *time.Time `json:"targetEndDate,omitempty"`
	Budget        float64    `json:"budget"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Total    int               `json:"total"`
}

type TeamMemberResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type ProjectTeamResponse struct {
	Users []TeamMemberResponse `json:"users"`
}

func toProjectResponse(p *models.Project) ProjectResponse {
	return ProjectResponse{
		ID:            p.ID.String(),
		Name:          p.Name,
		Description:   p.Description,
		Status:        string(p.Status),
		Priority:      p.Priority,
		HealthScore:   p.HealthScore,
		Progress:      p.Progress,
		StartDate:     p.StartDate,
		TargetEndDate: p.TargetEndDate,
		Budget:        p.Budget,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}
