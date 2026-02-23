package dto

import (
	"time"

	"github.com/google/uuid"
)

// ===== Common DTOs =====

// ResponseMeta represents metadata in API responses
type ResponseMeta struct {
	Page       *int       `json:"page,omitempty"`
	PerPage    *int       `json:"perPage,omitempty"`
	Total      *int       `json:"total,omitempty"`
	HasMore    *bool      `json:"hasMore,omitempty"`
	NextCursor *string    `json:"nextCursor,omitempty"`
	Timestamp  time.Time  `json:"timestamp"`
	RequestID  string     `json:"requestId"`
}

// ApiError represents an API error
type ApiError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ApiResponse represents the standard API response wrapper
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ApiError   `json:"error,omitempty"`
	Meta    ResponseMeta `json:"meta"`
}

// PaginationParams represents common pagination query parameters
type PaginationParams struct {
	Limit  int    `form:"limit,default=20" binding:"min=1,max=100"`
	Offset int    `form:"offset,default=0" binding:"min=0"`
	Cursor string `form:"cursor,omitempty"`
}

// CursorPaginationMeta represents metadata for cursor-based pagination
type CursorPaginationMeta struct {
	HasMore    bool   `json:"hasMore"`
	NextCursor string `json:"nextCursor,omitempty"`
	PerPage    int    `json:"perPage"`
}

// OffsetPaginationMeta represents metadata for offset-based pagination
type OffsetPaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// ValidationErrorDetail represents a single validation error
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorDetails represents validation error details
type ValidationErrorDetails struct {
	Errors map[string][]string `json:"errors"`
}

// RateLimitDetails represents rate limit error details
type RateLimitDetails struct {
	Limit       int       `json:"limit"`
	Remaining   int       `json:"remaining"`
	ResetAt     time.Time `json:"resetAt"`
	RetryAfter  int       `json:"retryAfter"`
}

// NewSuccessResponse creates a new successful API response
func NewSuccessResponse(data interface{}, meta ResponseMeta) ApiResponse {
	return ApiResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
}

// NewErrorResponse creates a new error API response
func NewErrorResponse(code, message string, details map[string]interface{}, requestID string) ApiResponse {
	return ApiResponse{
		Success: false,
		Error: &ApiError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: ResponseMeta{
			Timestamp: time.Now().UTC(),
			RequestID: requestID,
		},
	}
}

// WithRequestID sets the request ID in response meta
func (r ApiResponse) WithRequestID(requestID string) ApiResponse {
	r.Meta.RequestID = requestID
	r.Meta.Timestamp = time.Now().UTC()
	return r
}

// BatchOperationResult represents the result of a batch operation
type BatchOperationResult struct {
	ID      string     `json:"id"`
	Status  string     `json:"status"`
	Error   *ApiError  `json:"error,omitempty"`
}

// BatchOperationResponse represents the response for batch operations
type BatchOperationResponse struct {
	Processed int                    `json:"processed"`
	Succeeded int                    `json:"succeeded"`
	Failed    int                    `json:"failed"`
	Results   []BatchOperationResult `json:"results"`
}

// UserContext represents the authenticated user context
type UserContext struct {
	UserID         uuid.UUID `json:"userId"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Role           string    `json:"role"`
	Email          string    `json:"email"`
}

// IsAuthenticated checks if the user is authenticated
func (u UserContext) IsAuthenticated() bool {
	return u.UserID != uuid.Nil
}

// IsAdmin checks if the user is an admin
func (u UserContext) IsAdmin() bool {
	return u.Role == "admin"
}

// CanAccessOrganization checks if user can access the given organization
func (u UserContext) CanAccessOrganization(orgID uuid.UUID) bool {
	return u.OrganizationID == orgID || u.IsAdmin()
}

// SortParams represents sorting parameters
type SortParams struct {
	SortBy    string `form:"sortBy,default=created_at"`
	SortOrder string `form:"sortOrder,default=desc" binding:"oneof=asc desc"`
}

// DateRangeParams represents date range query parameters
type DateRangeParams struct {
	StartDate *time.Time `form:"startDate" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"endDate" time_format:"2006-01-02"`
}

// IDResponse represents a response containing just an ID
type IDResponse struct {
	ID string `json:"id"`
}

// StatusResponse represents a response containing just a status
type StatusResponse struct {
	Status string `json:"status"`
}

// CountResponse represents a response containing a count
type CountResponse struct {
	Count int `json:"count"`
}
