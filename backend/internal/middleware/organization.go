package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/SimpleAjax/Xephyr/internal/dto"
)

// OrganizationMiddleware validates organization access
type OrganizationMiddleware struct {
	// In real implementation, this would check organization membership
}

// NewOrganizationMiddleware creates a new organization middleware
func NewOrganizationMiddleware() *OrganizationMiddleware {
	return &OrganizationMiddleware{}
}

// RequireOrganization ensures an organization ID is provided and valid
func (m *OrganizationMiddleware) RequireOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orgID := ctx.GetHeader("X-Organization-Id")
		
		if orgID == "" {
			// Check if it's in the context (set by auth middleware)
			orgID = ctx.GetString("organizationId")
		}

		if orgID == "" {
			ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse(
				"VALIDATION_ERROR",
				"Organization ID is required",
				nil,
				GetRequestID(ctx),
			))
			ctx.Abort()
			return
		}

		// Validate UUID format
		if _, err := uuid.Parse(orgID); err != nil {
			ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse(
				"VALIDATION_ERROR",
				"Invalid organization ID format",
				nil,
				GetRequestID(ctx),
			))
			ctx.Abort()
			return
		}

		// TODO: Verify user has access to this organization
		// For now, we'll just accept it

		ctx.Set("organizationId", orgID)
		ctx.Next()
	}
}

// ExtractOrganization extracts organization ID if present but doesn't require it
func (m *OrganizationMiddleware) ExtractOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orgID := ctx.GetHeader("X-Organization-Id")
		if orgID != "" {
			if _, err := uuid.Parse(orgID); err == nil {
				ctx.Set("organizationId", orgID)
			}
		}
		ctx.Next()
	}
}
