package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// APIClient provides a high-level interface for E2E API testing
type APIClient struct {
	BaseURL        string
	APIToken       string
	OrganizationID string
	HTTPClient     *http.Client
	ctx            context.Context
}

// NewAPIClient creates a new API client for E2E tests
func NewAPIClient(baseURL, apiToken, orgID string) *APIClient {
	return &APIClient{
		BaseURL:        baseURL,
		APIToken:       apiToken,
		OrganizationID: orgID,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		ctx: context.Background(),
	}
}

// WithContext sets the context for the API client
func (c *APIClient) WithContext(ctx context.Context) *APIClient {
	c.ctx = ctx
	return c
}

// RequestOption is a functional option for customizing requests
type RequestOption func(*http.Request)

// WithHeader adds a custom header to the request
func WithHeader(key, value string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

// WithQueryParam adds a query parameter to the request
func WithQueryParam(key, value string) RequestOption {
	return func(req *http.Request) {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
	}
}

// Response is a generic API response wrapper
type Response[T any] struct {
	Success bool     `json:"success"`
	Data    T        `json:"data,omitempty"`
	Error   APIError `json:"error,omitempty"`
	Meta    Meta     `json:"meta,omitempty"`
}

// APIError represents an API error response
type APIError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string][]string `json:"details,omitempty"`
}

// Meta represents response metadata
type Meta struct {
	Page        int    `json:"page,omitempty"`
	PerPage     int    `json:"perPage,omitempty"`
	Total       int    `json:"total,omitempty"`
	Timestamp   string `json:"timestamp"`
	RequestID   string `json:"requestId"`
	HasMore     bool   `json:"hasMore,omitempty"`
	NextCursor  string `json:"nextCursor,omitempty"`
}

// Error implements the error interface
func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// IsSuccess returns true if the response indicates success
func (r Response[T]) IsSuccess() bool {
	return r.Success
}

// Get performs a GET request
func (c *APIClient) Get(path string, options ...RequestOption) (*http.Response, error) {
	return c.doRequest(http.MethodGet, path, nil, options...)
}

// GetJSON performs a GET request and parses the JSON response
func GetJSON[T any](c *APIClient, path string, options ...RequestOption) (*Response[T], error) {
	resp, err := c.Get(path, options...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseJSONResponse[T](resp)
}

// Post performs a POST request
func (c *APIClient) Post(path string, body interface{}, options ...RequestOption) (*http.Response, error) {
	return c.doRequest(http.MethodPost, path, body, options...)
}

// PostJSON performs a POST request and parses the JSON response
func PostJSON[T any](c *APIClient, path string, body interface{}, options ...RequestOption) (*Response[T], error) {
	resp, err := c.Post(path, body, options...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseJSONResponse[T](resp)
}

// Put performs a PUT request
func (c *APIClient) Put(path string, body interface{}, options ...RequestOption) (*http.Response, error) {
	return c.doRequest(http.MethodPut, path, body, options...)
}

// PutJSON performs a PUT request and parses the JSON response
func PutJSON[T any](c *APIClient, path string, body interface{}, options ...RequestOption) (*Response[T], error) {
	resp, err := c.Put(path, body, options...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseJSONResponse[T](resp)
}

// Patch performs a PATCH request
func (c *APIClient) Patch(path string, body interface{}, options ...RequestOption) (*http.Response, error) {
	return c.doRequest(http.MethodPatch, path, body, options...)
}

// PatchJSON performs a PATCH request and parses the JSON response
func PatchJSON[T any](c *APIClient, path string, body interface{}, options ...RequestOption) (*Response[T], error) {
	resp, err := c.Patch(path, body, options...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseJSONResponse[T](resp)
}

// Delete performs a DELETE request
func (c *APIClient) Delete(path string, options ...RequestOption) (*http.Response, error) {
	return c.doRequest(http.MethodDelete, path, nil, options...)
}

// DeleteJSON performs a DELETE request and parses the JSON response
func DeleteJSON[T any](c *APIClient, path string, options ...RequestOption) (*Response[T], error) {
	resp, err := c.Delete(path, options...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseJSONResponse[T](resp)
}

// doRequest performs the HTTP request
func (c *APIClient) doRequest(method, path string, body interface{}, options ...RequestOption) (*http.Response, error) {
	url, err := url.JoinPath(c.BaseURL, "/api/v1", path)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(c.ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set default headers
	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("X-Organization-Id", c.OrganizationID)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply custom options
	for _, opt := range options {
		opt(req)
	}

	return c.HTTPClient.Do(req)
}

// parseJSONResponse parses the HTTP response into a typed Response
func parseJSONResponse[T any](resp *http.Response) (*Response[T], error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result Response[T]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &result, nil
}

// ExpectStatus asserts the response has the expected status code
func ExpectStatus(resp *http.Response, expected int) error {
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("expected status %d, got %d: %s", expected, resp.StatusCode, string(body))
	}
	return nil
}

// IsSuccessStatus returns true for 2xx status codes
func IsSuccessStatus(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	ShouldRetry func(*http.Response, error) bool
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       1 * time.Second,
		ShouldRetry: func(resp *http.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500 || resp.StatusCode == 429
		},
	}
}

// DoWithRetry performs a request with retry logic
func (c *APIClient) DoWithRetry(method, path string, body interface{}, config RetryConfig, options ...RequestOption) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		resp, err = c.doRequest(method, path, body, options...)
		
		if !config.ShouldRetry(resp, err) {
			return resp, err
		}

		if attempt < config.MaxAttempts {
			if resp != nil {
				resp.Body.Close()
			}
			time.Sleep(config.Delay * time.Duration(attempt))
		}
	}

	return resp, err
}
