package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", "X-Organization-Id"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// CORS middleware handles Cross-Origin Resource Sharing
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, o := range config.AllowOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*" {
				ctx.Header("Access-Control-Allow-Origin", "*")
			} else {
				ctx.Header("Access-Control-Allow-Origin", origin)
			}
		}

		ctx.Header("Access-Control-Allow-Methods", joinStrings(config.AllowMethods, ", "))
		ctx.Header("Access-Control-Allow-Headers", joinStrings(config.AllowHeaders, ", "))
		ctx.Header("Access-Control-Expose-Headers", joinStrings(config.ExposeHeaders, ", "))
		ctx.Header("Access-Control-Allow-Credentials", boolToString(config.AllowCredentials))
		ctx.Header("Access-Control-Max-Age", intToString(config.MaxAge))

		// Handle preflight requests
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

// Helper functions
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func intToString(i int) string {
	return string(rune(i))
}
