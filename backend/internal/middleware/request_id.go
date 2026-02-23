package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

// RequestID middleware generates or extracts request ID
func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Set("requestId", requestID)
		ctx.Header(RequestIDHeader, requestID)

		ctx.Next()
	}
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx *gin.Context) string {
	return ctx.GetString("requestId")
}
