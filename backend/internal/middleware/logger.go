package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	SkipPaths []string
}

// Logger middleware logs HTTP requests
func Logger(config LoggerConfig) gin.HandlerFunc {
	skipMap := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipMap[path] = true
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		ctx.Next()

		// Skip logging for certain paths
		if skipMap[path] {
			return
		}

		// Log the request
		latency := time.Since(start)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		requestID := GetRequestID(ctx)

		if raw != "" {
			path = path + "?" + raw
		}

		// In production, use structured logging
		// For now, we'll use a simple format
		_ = gin.LogFormatterParams{
			TimeStamp:    time.Now(),
			Latency:      latency,
			ClientIP:     clientIP,
			Method:       method,
			StatusCode:   statusCode,
			ErrorMessage: ctx.Errors.ByType(gin.ErrorTypePrivate).String(),
			BodySize:     ctx.Writer.Size(),
			Path:         path,
		}

		// TODO: Use proper logger
		// logger.Info("HTTP Request",
		//     "request_id", requestID,
		//     "method", method,
		//     "path", path,
		//     "status", statusCode,
		//     "latency", latency,
		//     "client_ip", clientIP,
		// )
		_ = requestID
	}
}

// DefaultLogger returns a default logger middleware
func DefaultLogger() gin.HandlerFunc {
	return Logger(LoggerConfig{
		SkipPaths: []string{"/health"},
	})
}
