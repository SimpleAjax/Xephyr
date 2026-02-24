package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SimpleAjax/Xephyr/internal/routes"
	"github.com/SimpleAjax/Xephyr/tests/fixtures"
)

// E2E test configuration
type E2EConfig struct {
	BaseURL          string
	APIToken         string
	OrganizationID   string
	TestTimeout      time.Duration
	SkipCleanup      bool
	GenerateTestData bool
}

var (
	// Global test configuration
	Config E2EConfig

	// Test server instance
	TestServer *httptest.Server

	// HTTP client for tests
	HTTPClient *http.Client

	// Context for test lifecycle management
	TestCtx   context.Context
	CancelCtx context.CancelFunc
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)

	// Run specs
	RunSpecs(t, "Xephyr E2E Test Suite")
}

var _ = BeforeSuite(func() {
	// Initialize test context
	TestCtx, CancelCtx = context.WithTimeout(context.Background(), 30*time.Minute)

	// Load configuration from environment or use defaults
	Config = E2EConfig{
		BaseURL:          getEnv("E2E_BASE_URL", ""),
		APIToken:         getEnv("E2E_API_TOKEN", "test-token"),
		OrganizationID:   getEnv("E2E_ORG_ID", "550e8400-e29b-41d4-a716-446655440000"),
		TestTimeout:      60 * time.Second,
		SkipCleanup:      getEnv("E2E_SKIP_CLEANUP", "false") == "true",
		GenerateTestData: true,
	}

	// Setup test HTTP client
	HTTPClient = &http.Client{
		Timeout: Config.TestTimeout,
	}

	// If no external URL provided, start test server
	if Config.BaseURL == "" {
		TestServer = setupTestServer()
		Config.BaseURL = TestServer.URL
	}

	fmt.Printf("\n🚀 E2E Test Suite Started\n")
	fmt.Printf("   Base URL: %s\n", Config.BaseURL)
	fmt.Printf("   Organization: %s\n", Config.OrganizationID)
	fmt.Println()

	// Generate test fixtures if needed
	if Config.GenerateTestData {
		generateTestData()
	}
})

var _ = AfterSuite(func() {
	// Cleanup test data unless skipped
	if !Config.SkipCleanup {
		cleanupTestData()
	}

	// Shutdown test server
	if TestServer != nil {
		TestServer.Close()
	}

	// Cancel context
	if CancelCtx != nil {
		CancelCtx()
	}

	fmt.Printf("\n✅ E2E Test Suite Completed\n")
})

// setupTestServer creates a test HTTP server with all routes
func setupTestServer() *httptest.Server {
	router := routes.SetupRoutes().GetEngine()
	return httptest.NewServer(router)
}

// generateTestData creates test fixtures in the database
func generateTestData() {
	// Create test teams, projects, tasks
	_ = fixtures.CreateTestTeam()
	_ = fixtures.CreateTestSkills()
	_ = fixtures.CreateTestProjects()
	_ = fixtures.CreateTestTasks()
	_ = fixtures.CreateTestWorkloadData()
}

// cleanupTestData removes test data from the database
func cleanupTestData() {
	// Cleanup logic would go here
	// This would remove all test-generated data
}

// getEnv retrieves environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// APIRequest represents an API request builder
type APIRequest struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
	Query   map[string]string
}

// APIResponse represents a typed API response
type APIResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// ExecuteRequest performs an HTTP request against the test server
func ExecuteRequest(req APIRequest) (*APIResponse, error) {
	url := Config.BaseURL + "/api/v1" + req.Path

	// Build request
	httpReq, err := buildHTTPRequest(req.Method, url, req.Body)
	if err != nil {
		return nil, err
	}

	// Add headers
	httpReq.Header.Set("Authorization", "Bearer "+Config.APIToken)
	httpReq.Header.Set("X-Organization-Id", Config.OrganizationID)
	httpReq.Header.Set("Content-Type", "application/json")

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Add query parameters
	if len(req.Query) > 0 {
		q := httpReq.URL.Query()
		for key, value := range req.Query {
			q.Add(key, value)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	// Execute request
	resp, err := HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &APIResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
	}, nil
}

// buildHTTPRequest creates an HTTP request with optional body
func buildHTTPRequest(method, url string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}
	return http.NewRequestWithContext(TestCtx, method, url, bodyReader)
}

// ParseResponse unmarshals API response into target struct
func ParseResponse(resp *APIResponse, target interface{}) error {
	return json.Unmarshal(resp.Body, target)
}
