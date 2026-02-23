package backend_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/xephyr-ai/xephyr-backend/test/reporters"
)

func TestXephyrBackend(t *testing.T) {
	RegisterFailHandler(Fail)
	
	// Create custom reporter for Gherkin-style output
	gherkinReporter := reporters.NewGherkinReporter()
	
	// Run specs with custom reporter
	RunSpecs(t, "Xephyr Backend Test Suite", Label("xephyr", "backend"), ReportAfterSuite("", func(report Report) {
		// Generate JSON output for UI consumption
		jsonOutput, err := gherkinReporter.GenerateJSONOutput()
		if err == nil {
			// In real implementation, save to file or send to reporting service
			_ = jsonOutput
		}
		
		// Print final summary
		gherkinReporter.PrintSummary()
	}))
}

// TestConfiguration provides global test configuration
type TestConfiguration struct {
	// Database configuration for integration tests
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	
	// Test behavior flags
	SkipIntegrationTests bool
	GenerateCoverage     bool
}

// GlobalConfig is the global test configuration
var GlobalConfig TestConfiguration

func init() {
	// Initialize default configuration
	GlobalConfig = TestConfiguration{
		DBHost:               "localhost",
		DBPort:               "5432",
		DBName:               "xephyr_test",
		DBUser:               "postgres",
		DBPassword:           "postgres",
		SkipIntegrationTests: false,
		GenerateCoverage:     true,
	}
}

// BeforeSuite runs before all tests
var _ = BeforeSuite(func() {
	// Global setup
	// - Initialize test database
	// - Set up test fixtures
	// - Configure logging
})

// AfterSuite runs after all tests
var _ = AfterSuite(func() {
	// Global teardown
	// - Clean up test database
	// - Close connections
})
