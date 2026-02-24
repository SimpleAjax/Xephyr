package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/SimpleAjax/Xephyr/internal/dto"
)

// AuthMiddleware validates JWT tokens and sets user context
type AuthMiddleware struct {
	// In real implementation, this would have a JWT validator
	jwtSecret string
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

// Authenticate validates the JWT token and sets user context
// For development: if no auth header is provided, generates dummy user/org IDs
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		
		// TODO: Implement actual JWT validation
		// For development/demo: allow requests without auth header
		userID := ctx.GetHeader("X-User-ID")
		orgID := ctx.GetHeader("X-Organization-Id")

		// If auth header provided, extract token (but don't validate for now)
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				// Token provided - could validate here in production
			}
		}

		if userID == "" {
			// Generate dummy user ID for demo
			userID = uuid.New().String()
		}
		if orgID == "" {
			// Generate dummy org ID for demo
			orgID = uuid.New().String()
		}

		// Validate UUIDs
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dto.NewErrorResponse(
				"UNAUTHORIZED",
				"Invalid user ID",
				nil,
				GetRequestID(ctx),
			))
			ctx.Abort()
			return
		}

		orgUUID, err := uuid.Parse(orgID)
		if err != nil {
			ctx.JSON(http.StatusForbidden, dto.NewErrorResponse(
				"FORBIDDEN",
				"Invalid organization ID",
				nil,
				GetRequestID(ctx),
			))
			ctx.Abort()
			return
		}

		// Set user context
		ctx.Set("userId", userUUID.String())
		ctx.Set("organizationId", orgUUID.String())
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				ctx.Set("token", parts[1])
			}
		}

		ctx.Next()
	}
}

// OptionalAuth allows requests without authentication but sets user info if available
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" {
			// Try to extract and validate token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				userID := ctx.GetHeader("X-User-ID")
				orgID := ctx.GetHeader("X-Organization-Id")

				if userID != "" && orgID != "" {
					ctx.Set("userId", userID)
					ctx.Set("organizationId", orgID)
				}
			}
		}
		ctx.Next()
	}
}

// RequireRole checks if the user has the required role
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: Implement role checking
		// For now, just pass through
		ctx.Next()
	}
}

// GetUserID retrieves the user ID from context
func GetUserID(ctx *gin.Context) string {
	return ctx.GetString("userId")
}

// GetOrganizationID retrieves the organization ID from context
func GetOrganizationID(ctx *gin.Context) string {
	return ctx.GetString("organizationId")
}
