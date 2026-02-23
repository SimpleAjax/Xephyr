package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xephyr-ai/xephyr-backend/internal/dto"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
	mu       sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimit middleware limits requests per client
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Use IP + UserID as key if authenticated, otherwise just IP
		key := ctx.ClientIP()
		if userID := GetUserID(ctx); userID != "" {
			key = userID
		}

		allowed, remaining, resetAt := rl.allowRequest(key)

		// Set rate limit headers
		ctx.Header("X-RateLimit-Limit", intToString(rl.limit))
		ctx.Header("X-RateLimit-Remaining", intToString(remaining))
		ctx.Header("X-RateLimit-Reset", resetAt.Format(time.RFC3339))

		if !allowed {
			ctx.JSON(http.StatusTooManyRequests, dto.NewErrorResponse(
				"RATE_LIMITED",
				"Rate limit exceeded. Please try again later.",
				map[string]interface{}{
					"limit":       rl.limit,
					"remaining":   0,
					"resetAt":     resetAt,
					"retryAfter":  int(time.Until(resetAt).Seconds()),
				},
				GetRequestID(ctx),
			))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// allowRequest checks if a request is allowed and returns rate limit info
func (rl *RateLimiter) allowRequest(key string) (bool, int, time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Clean old requests
	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, req := range requests {
			if req.After(windowStart) {
				validRequests = append(validRequests, req)
			}
		}
		rl.requests[key] = validRequests
	}

	// Check if limit is exceeded
	requests := rl.requests[key]
	if len(requests) >= rl.limit {
		resetAt := requests[0].Add(rl.window)
		return false, 0, resetAt
	}

	// Add current request
	rl.requests[key] = append(requests, now)
	remaining := rl.limit - len(rl.requests[key])
	resetAt := now.Add(rl.window)

	return true, remaining, resetAt
}

// Cleanup removes old entries (should be called periodically)
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	windowStart := time.Now().Add(-rl.window)
	for key, requests := range rl.requests {
		var validRequests []time.Time
		for _, req := range requests {
			if req.After(windowStart) {
				validRequests = append(validRequests, req)
			}
		}
		if len(validRequests) == 0 {
			delete(rl.requests, key)
		} else {
			rl.requests[key] = validRequests
		}
	}
}
