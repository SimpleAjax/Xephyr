package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/internal/repositories"
)

// TaskController handles task CRUD HTTP requests
type TaskController struct {
	repos repositories.Repositories
}

// NewTaskController creates a new task controller
func NewTaskController(repos repositories.Repositories) *TaskController {
	return &TaskController{repos: repos}
}

// ListTasks godoc
// @Summary List all tasks
// @Description Get a list of tasks, optionally filtered by project
// @Tags tasks
// @Accept json
// @Produce json
// @Param projectId query string false "Filter by project ID"
// @Param status query string false "Filter by status"
// @Param assigneeId query string false "Filter by assignee"
// @Param limit query int false "Limit results" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} dto.ApiResponse{data=TaskListResponse}
// @Security BearerAuth
// @Router /tasks [get]
func (c *TaskController) ListTasks(ctx *gin.Context) {
	projectID := ctx.Query("projectId")
	status := ctx.Query("status")
	assigneeID := ctx.Query("assigneeId")

	var params repositories.ListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		params = repositories.ListParams{Limit: 50, Offset: 0}
	}
	if params.Limit == 0 {
		params.Limit = 50
	}

	var tasks []models.Task
	var total int64
	var err error

	if projectID != "" {
		projUUID, err := uuid.Parse(projectID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid project ID", nil, ctx.GetString("requestId")))
			return
		}
		tasks, total, err = c.repos.GetTask().ListByProject(ctx.Request.Context(), projUUID, params)
	} else if assigneeID != "" {
		userUUID, err := uuid.Parse(assigneeID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid assignee ID", nil, ctx.GetString("requestId")))
			return
		}
		tasks, total, err = c.repos.GetTask().ListByAssignee(ctx.Request.Context(), userUUID, params)
	} else if status != "" {
		taskStatus := models.TaskStatus(status)
		tasks, total, err = c.repos.GetTask().ListByStatus(ctx.Request.Context(), taskStatus, params)
	} else {
		tasks, total, err = c.repos.GetTask().List(ctx.Request.Context(), params)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	response := TaskListResponse{
		Tasks: make([]TaskResponse, 0, len(tasks)),
		Total: int(total),
	}

	for _, t := range tasks {
		response.Tasks = append(response.Tasks, toTaskResponse(&t))
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(response, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetTask godoc
// @Summary Get a task by ID
// @Description Get detailed information about a specific task
// @Tags tasks
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} dto.ApiResponse{data=TaskResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /tasks/{taskId} [get]
func (c *TaskController) GetTask(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid task ID", nil, ctx.GetString("requestId")))
		return
	}

	task, err := c.repos.GetTask().GetByID(ctx.Request.Context(), taskUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Task not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(toTaskResponse(task), dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task in a project
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Task creation request"
// @Success 201 {object} dto.ApiResponse{data=TaskResponse}
// @Failure 400 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /tasks [post]
func (c *TaskController) CreateTask(ctx *gin.Context) {
	var req CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	projUUID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid project ID", nil, ctx.GetString("requestId")))
		return
	}

	task := &models.Task{
		ProjectID:      projUUID,
		Title:          req.Title,
		Description:    req.Description,
		Status:         models.TaskStatusBacklog,
		Priority:       models.TaskPriority(req.Priority),
		EstimatedHours: req.EstimatedHours,
		HierarchyLevel: req.HierarchyLevel,
		IsMilestone:    req.IsMilestone,
	}

	if req.Status != "" {
		task.Status = models.TaskStatus(req.Status)
	}
	if req.ParentTaskID != "" {
		parentUUID, err := uuid.Parse(req.ParentTaskID)
		if err == nil {
			task.ParentTaskID = &parentUUID
		}
	}
	if req.AssigneeID != "" {
		assigneeUUID, err := uuid.Parse(req.AssigneeID)
		if err == nil {
			task.AssigneeID = &assigneeUUID
		}
	}

	if err := c.repos.GetTask().Create(ctx.Request.Context(), task); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusCreated, dto.NewSuccessResponse(toTaskResponse(task), dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Param request body UpdateTaskRequest true "Task update request"
// @Success 200 {object} dto.ApiResponse{data=TaskResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /tasks/{taskId} [patch]
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid task ID", nil, ctx.GetString("requestId")))
		return
	}

	task, err := c.repos.GetTask().GetByID(ctx.Request.Context(), taskUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "Task not found", nil, ctx.GetString("requestId")))
		return
	}

	var req UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	// Update fields
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = models.TaskStatus(req.Status)
	}
	if req.Priority != "" {
		task.Priority = models.TaskPriority(req.Priority)
	}
	if req.EstimatedHours != 0 {
		task.EstimatedHours = req.EstimatedHours
	}
	if req.ActualHours != 0 {
		task.ActualHours = req.ActualHours
	}
	if req.AssigneeID != nil {
		if *req.AssigneeID == "" {
			task.AssigneeID = nil
		} else {
			assigneeUUID, err := uuid.Parse(*req.AssigneeID)
			if err == nil {
				task.AssigneeID = &assigneeUUID
			}
		}
	}

	if err := c.repos.GetTask().Update(ctx.Request.Context(), task); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(toTaskResponse(task), dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// UpdateTaskStatus godoc
// @Summary Update task status
// @Description Update just the status of a task
// @Tags tasks
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Param request body UpdateTaskStatusRequest true "Status update request"
// @Success 200 {object} dto.ApiResponse{data=TaskResponse}
// @Security BearerAuth
// @Router /tasks/{taskId}/status [post]
func (c *TaskController) UpdateTaskStatus(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid task ID", nil, ctx.GetString("requestId")))
		return
	}

	var req UpdateTaskStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	taskStatus := models.TaskStatus(req.Status)
	if err := c.repos.GetTask().UpdateStatus(ctx.Request.Context(), taskUUID, taskStatus); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	// Fetch updated task
	task, err := c.repos.GetTask().GetByID(ctx.Request.Context(), taskUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(toTaskResponse(task), dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// AssignTask godoc
// @Summary Assign task to user
// @Description Assign a task to a specific user
// @Tags tasks
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Param request body AssignTaskRequest true "Assign request"
// @Success 200 {object} dto.ApiResponse{data=TaskResponse}
// @Security BearerAuth
// @Router /tasks/{taskId}/assign [post]
func (c *TaskController) AssignTask(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid task ID", nil, ctx.GetString("requestId")))
		return
	}

	var req AssignTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	var assigneeID *uuid.UUID
	if req.PersonID != "" {
		userUUID, err := uuid.Parse(req.PersonID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid user ID", nil, ctx.GetString("requestId")))
			return
		}
		assigneeID = &userUUID
	}

	if err := c.repos.GetTask().UpdateAssignee(ctx.Request.Context(), taskUUID, assigneeID); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	task, err := c.repos.GetTask().GetByID(ctx.Request.Context(), taskUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(toTaskResponse(task), dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Soft-delete a task
// @Tags tasks
// @Accept json
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 204
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /tasks/{taskId} [delete]
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid task ID", nil, ctx.GetString("requestId")))
		return
	}

	if err := c.repos.GetTask().Delete(ctx.Request.Context(), taskUUID); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Request/Response types

type CreateTaskRequest struct {
	ProjectID      string  `json:"projectId" binding:"required"`
	ParentTaskID   string  `json:"parentTaskId,omitempty"`
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description"`
	Status         string  `json:"status,omitempty"`
	Priority       string  `json:"priority,omitempty"`
	EstimatedHours float64 `json:"estimatedHours"`
	HierarchyLevel int     `json:"hierarchyLevel,omitempty"`
	AssigneeID     string  `json:"assigneeId,omitempty"`
	IsMilestone    bool    `json:"isMilestone,omitempty"`
}

type UpdateTaskRequest struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Status         string  `json:"status"`
	Priority       string  `json:"priority"`
	EstimatedHours float64 `json:"estimatedHours"`
	ActualHours    float64 `json:"actualHours"`
	AssigneeID     *string `json:"assigneeId"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AssignTaskRequest struct {
	PersonID string `json:"personId"`
	Note     string `json:"note"`
}

type TaskResponse struct {
	ID             string     `json:"id"`
	ProjectID      string     `json:"projectId"`
	ParentTaskID   *string    `json:"parentTaskId,omitempty"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	EstimatedHours float64    `json:"estimatedHours"`
	ActualHours    float64    `json:"actualHours"`
	AssigneeID     *string    `json:"assigneeId,omitempty"`
	AssigneeName   string     `json:"assigneeName,omitempty"`
	IsMilestone    bool       `json:"isMilestone"`
	IsCriticalPath bool       `json:"isCriticalPath"`
	RiskScore      int        `json:"riskScore"`
	DueDate        *time.Time `json:"dueDate,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
	Total int            `json:"total"`
}

func toTaskResponse(t *models.Task) TaskResponse {
	resp := TaskResponse{
		ID:             t.ID.String(),
		ProjectID:      t.ProjectID.String(),
		Title:          t.Title,
		Description:    t.Description,
		Status:         string(t.Status),
		Priority:       string(t.Priority),
		EstimatedHours: t.EstimatedHours,
		ActualHours:    t.ActualHours,
		IsMilestone:    t.IsMilestone,
		IsCriticalPath: t.IsCriticalPath,
		RiskScore:      t.RiskScore,
		DueDate:        t.DueDate,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}

	if t.ParentTaskID != nil {
		parentID := t.ParentTaskID.String()
		resp.ParentTaskID = &parentID
	}

	if t.AssigneeID != nil {
		assigneeID := t.AssigneeID.String()
		resp.AssigneeID = &assigneeID
		if t.Assignee != nil {
			resp.AssigneeName = t.Assignee.Name
		}
	}

	return resp
}
