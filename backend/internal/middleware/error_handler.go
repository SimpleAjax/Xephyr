package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
)

// ErrorHandler middleware handles panics and returns proper error responses
func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error (in production, use proper logging)
				// logger.Error("Panic recovered", err)

				ctx.JSON(http.StatusInternalServerError, dto.NewErrorResponse(
					"INTERNAL_ERROR",
					"An unexpected error occurred",
					nil,
					GetRequestID(ctx),
				))
				ctx.Abort()
			}
		}()

		ctx.Next()

		// Handle errors set by gin
		if len(ctx.Errors) > 0 {
			lastError := ctx.Errors.Last()
			
			var statusCode int
			switch lastError.Type {
			case gin.ErrorTypeBind:
				statusCode = http.StatusBadRequest
			case gin.ErrorTypeRender:
				statusCode = http.StatusUnprocessableEntity
			default:
				statusCode = http.StatusInternalServerError
			}

			ctx.JSON(statusCode, dto.NewErrorResponse(
				"REQUEST_ERROR",
				lastError.Error(),
				nil,
				GetRequestID(ctx),
			))
		}
	}
}

// NotFoundHandler handles 404 errors
func NotFoundHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.NewErrorResponse(
			"NOT_FOUND",
			"The requested resource was not found",
			nil,
			GetRequestID(ctx),
		))
	}
}

// MethodNotAllowedHandler handles 405 errors
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, dto.NewErrorResponse(
			"METHOD_NOT_ALLOWED",
			"The requested method is not allowed for this resource",
			nil,
			GetRequestID(ctx),
		))
	}
}
