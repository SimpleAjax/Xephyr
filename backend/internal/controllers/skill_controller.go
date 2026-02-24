package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/SimpleAjax/Xephyr/internal/dto"
	"github.com/SimpleAjax/Xephyr/internal/models"
	"github.com/SimpleAjax/Xephyr/internal/repositories"
)

// SkillController handles skill HTTP requests
type SkillController struct {
	repos repositories.Repositories
}

// NewSkillController creates a new skill controller
func NewSkillController(repos repositories.Repositories) *SkillController {
	return &SkillController{repos: repos}
}

// ListSkills godoc
// @Summary List all skills
// @Description Get a list of all skills (global and organization-specific)
// @Tags skills
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=SkillListResponse}
// @Security BearerAuth
// @Router /skills [get]
func (c *SkillController) ListSkills(ctx *gin.Context) {
	// For now, return skills from users in the organization
	orgID := ctx.GetString("organizationId")
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid organization ID", nil, ctx.GetString("requestId")))
		return
	}

	// Get users with their skills
	users, _, err := c.repos.GetUser().ListByOrganization(ctx.Request.Context(), orgUUID, repositories.ListParams{Limit: 100})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	// Collect unique skills
	skillMap := make(map[string]SkillResponse)
	for _, user := range users {
		// Need to preload skills - for now, fetch each user individually
		userWithSkills, err := c.repos.GetUser().GetByID(ctx.Request.Context(), user.ID)
		if err != nil {
			continue
		}
		for _, s := range userWithSkills.Skills {
			skillMap[s.Skill.ID.String()] = SkillResponse{
				ID:          s.Skill.ID.String(),
				Name:        s.Skill.Name,
				Category:    s.Skill.Category,
				Description: s.Skill.Description,
			}
		}
	}

	skills := make([]SkillResponse, 0, len(skillMap))
	for _, s := range skillMap {
		skills = append(skills, s)
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(SkillListResponse{
		Skills: skills,
		Total:  len(skills),
	}, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// GetSkillGaps godoc
// @Summary Get skill gaps
// @Description Get skills required by tasks but not available in the team
// @Tags skills
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=SkillGapsResponse}
// @Security BearerAuth
// @Router /skills/gaps [get]
func (c *SkillController) GetSkillGaps(ctx *gin.Context) {
	orgID := ctx.GetString("organizationId")
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid organization ID", nil, ctx.GetString("requestId")))
		return
	}

	// Get all projects in organization
	projects, _, err := c.repos.GetProject().ListByOrganization(ctx.Request.Context(), orgUUID, repositories.ListParams{Limit: 100})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_ERROR", err.Error(), nil, ctx.GetString("requestId")))
		return
	}

	// Collect skill gaps (simplified implementation)
	gaps := []SkillGap{}

	// Check for nudges with skill_gap type
	nudgeFilters := repositories.NudgeFilters{}
	skillGapType := models.NudgeTypeSkillGap
	nudgeFilters.Type = &skillGapType

	nudges, _, err := c.repos.GetNudge().List(ctx.Request.Context(), orgUUID, nudgeFilters, repositories.ListParams{Limit: 50})
	if err == nil {
		for _, nudge := range nudges {
			if nudge.Metrics != nil {
				if skillName, ok := nudge.Metrics["requiredSkill"].(string); ok {
					gaps = append(gaps, SkillGap{
						SkillID:   "",
						SkillName: skillName,
						Reason:    nudge.Description,
					})
				}
			}
		}
	}

	// For each project, check task skills
	for _, project := range projects {
		tasks, _, err := c.repos.GetTask().ListByProject(ctx.Request.Context(), project.ID, repositories.ListParams{Limit: 100})
		if err != nil {
			continue
		}

		for _, task := range tasks {
			// Get task with skills
			taskWithSkills, err := c.repos.GetTask().GetByID(ctx.Request.Context(), task.ID)
			if err != nil {
				continue
			}

			for _, taskSkill := range taskWithSkills.Skills {
				// Check if any user has this skill
				hasSkill := false
				users, _, _ := c.repos.GetUser().ListByOrganization(ctx.Request.Context(), orgUUID, repositories.ListParams{Limit: 100})
				for _, user := range users {
					userWithSkills, err := c.repos.GetUser().GetByID(ctx.Request.Context(), user.ID)
					if err != nil {
						continue
					}
					for _, userSkill := range userWithSkills.Skills {
						if userSkill.SkillID == taskSkill.SkillID {
							hasSkill = true
							break
						}
					}
					if hasSkill {
						break
					}
				}

				if !hasSkill {
					gaps = append(gaps, SkillGap{
						SkillID:   taskSkill.SkillID.String(),
						SkillName: taskSkill.Skill.Name,
						Reason:    "Required by task: " + task.Title,
					})
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(SkillGapsResponse{
		Gaps: gaps,
	}, dto.ResponseMeta{
		Timestamp: getTimestamp(),
		RequestID: ctx.GetString("requestId"),
	}))
}

// Request/Response types

type SkillResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type SkillListResponse struct {
	Skills []SkillResponse `json:"skills"`
	Total  int             `json:"total"`
}

type SkillGap struct {
	SkillID   string `json:"skillId"`
	SkillName string `json:"skillName"`
	Reason    string `json:"reason"`
}

type SkillGapsResponse struct {
	Gaps []SkillGap `json:"gaps"`
}
