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

// UserController handles user CRUD HTTP requests
type UserController struct {
	repos repositories.Repositories
}

// NewUserController creates a new user controller
func NewUserController(repos repositories.Repositories) *UserController {
	return &UserController{repos: repos}
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of users in the organization
// @Tags users
// @Accept json
// @Produce json
// @Param role query string false "Filter by role"
// @Param limit query int false "Limit results" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} dto.ApiResponse{data=UserListResponse}
// @Security BearerAuth
// @Router /users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	orgID := ctx.GetString("organizationId")
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid organization ID", nil, ctx.GetString("requestId")))
		return
	}

	role := ctx.Query("role")

	var params repositories.ListParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		params = repositories.DefaultListParams()
	}

	var users []models.User
	var total int64

	if role != "" {
		userRole := models.UserRole(role)
		users, err = c.repos.GetUser().GetByOrganizationAndRole(ctx.Request.Context(), orgUUID, userRole)
		total = int64(len(users))
	} else {
		users, total, err = c.repos.GetUser().ListByOrganization(ctx.Request.Context(), orgUUID, params)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	response := UserListResponse{
		Users: make([]UserResponse, 0, len(users)),
		Total: int(total),
	}

	for _, u := range users {
		response.Users = append(response.Users, toUserResponse(&u))
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(response, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get detailed information about a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} dto.ApiResponse{data=UserResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /users/{userId} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("userId")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid user ID", nil, ctx.GetString("requestId")))
		return
	}

	user, err := c.repos.GetUser().GetByID(ctx.Request.Context(), userUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "User not found", nil, ctx.GetString("requestId")))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(toUserResponse(user), dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetUserSkills godoc
// @Summary Get user skills
// @Description Get skills for a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} dto.ApiResponse{data=UserSkillsResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /users/{userId}/skills [get]
func (c *UserController) GetUserSkills(ctx *gin.Context) {
	userID := ctx.Param("userId")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid user ID", nil, ctx.GetString("requestId")))
		return
	}

	user, err := c.repos.GetUser().GetByID(ctx.Request.Context(), userUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "User not found", nil, ctx.GetString("requestId")))
		return
	}

	skills := make([]UserSkillResponse, 0, len(user.Skills))
	for _, s := range user.Skills {
		skills = append(skills, UserSkillResponse{
			Skill: UserSkillDetail{
				ID:       s.Skill.ID.String(),
				Name:     s.Skill.Name,
				Category: s.Skill.Category,
			},
			Proficiency:       s.Proficiency,
			YearsOfExperience: s.YearsOfExperience,
		})
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(UserSkillsResponse{
		UserID: user.ID.String(),
		Skills: skills,
	}, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetUserWorkload godoc
// @Summary Get user workload
// @Description Get current workload for a specific user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} dto.ApiResponse{data=UserWorkloadResponse}
// @Failure 404 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /users/{userId}/workload [get]
func (c *UserController) GetUserWorkload(ctx *gin.Context) {
	userID := ctx.Param("userId")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid user ID", nil, ctx.GetString("requestId")))
		return
	}

	user, err := c.repos.GetUser().GetByID(ctx.Request.Context(), userUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse("NOT_FOUND", "User not found", nil, ctx.GetString("requestId")))
		return
	}

	// Get workload entries for current week
	weekStart := getCurrentWeekStart()
	entry, err := c.repos.GetWorkload().GetByUserAndWeek(ctx.Request.Context(), userUUID, weekStart)
	if err != nil {
		// Return empty workload
		ctx.JSON(http.StatusOK, dto.NewSuccessResponse(UserWorkloadResponse{
			UserID:              user.ID.String(),
			Name:                user.Name,
			AllocationPercentage: 0,
			AssignedHours:       0,
			CapacityHours:       40,
			AssignedTasks:       0,
		}, dto.ResponseMeta{
			Timestamp: getTimestamp(),
			RequestID: ctx.GetString("requestId"),
		}))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(UserWorkloadResponse{
		UserID:              user.ID.String(),
		Name:                user.Name,
		AllocationPercentage: entry.AllocationPercentage,
		AssignedHours:       entry.TotalEstimatedHours,
		CapacityHours:       entry.AvailableHours,
		AssignedTasks:       entry.AssignedTasks,
	}, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// Request/Response types

type UserResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	AvatarURL   string    `json:"avatarUrl"`
	HourlyRate  float64   `json:"hourlyRate"`
	Timezone    string    `json:"timezone"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}

type UserSkillResponse struct {
	Skill             UserSkillDetail `json:"skill"`
	Proficiency       int             `json:"proficiency"`
	YearsOfExperience float64         `json:"yearsOfExperience"`
}

type UserSkillDetail struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type UserSkillsResponse struct {
	UserID string              `json:"userId"`
	Skills []UserSkillResponse `json:"skills"`
}

type UserWorkloadResponse struct {
	UserID              string  `json:"userId"`
	Name                string  `json:"name"`
	AllocationPercentage int    `json:"allocationPercentage"`
	AssignedHours       float64 `json:"assignedHours"`
	CapacityHours       float64 `json:"capacityHours"`
	AssignedTasks       int     `json:"assignedTasks"`
}

func toUserResponse(u *models.User) UserResponse {
	return UserResponse{
		ID:         u.ID.String(),
		Email:      u.Email,
		Name:       u.Name,
		AvatarURL:    u.AvatarURL,
		HourlyRate: u.HourlyRate,
		Timezone:   u.Timezone,
		IsActive:   u.IsActive,
		CreatedAt:  u.CreatedAt,
	}
}

func getCurrentWeekStart() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
}
