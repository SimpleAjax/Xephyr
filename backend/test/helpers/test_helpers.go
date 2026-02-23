package helpers

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// RunSuite runs the test suite
func RunSuite(t *testing.T, suiteName string) {
	RegisterFailHandler(Fail)
	RunSpecs(t, suiteName)
}

// TestContext holds common test dependencies
type TestContext struct {
	// Add common dependencies here (DB, services, etc.)
}

// NewTestContext creates a new test context
func NewTestContext() *TestContext {
	return &TestContext{}
}

// Cleanup performs cleanup after tests
func (ctx *TestContext) Cleanup() {
	// Cleanup resources
}

// Ptr returns a pointer to the given value
func Ptr[T any](v T) *T {
	return &v
}

// MustUUID parses a UUID string and panics if invalid
func MustUUID(s string) string {
	// In a real implementation, this would validate and return a proper UUID
	// For now, we just return the string as-is assuming it's a valid UUID format
	return s
}

// IntPtr returns a pointer to an int
func IntPtr(i int) *int {
	return &i
}

// Float64Ptr returns a pointer to a float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	return &s
}
