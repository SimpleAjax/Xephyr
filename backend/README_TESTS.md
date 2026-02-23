# Xephyr Backend Test Suite

This directory contains the comprehensive TDD test suite for the Xephyr AI Project Management Platform backend.

## Test Structure

```
backend/
├── internal/
│   ├── models/              # Domain models (no tests here, tested via services)
│   ├── services/            # Business logic services with corresponding test files
│   │   ├── priority_service_test.go      # Priority Engine tests
│   │   ├── health_service_test.go        # Health Scoring tests
│   │   ├── nudge_service_test.go         # Nudge Detection tests
│   │   ├── assignment_service_test.go    # Assignment Engine tests
│   │   ├── dependency_service_test.go    # Dependency Management tests
│   │   ├── scenario_service_test.go      # Scenario Processing tests
│   │   └── progress_service_test.go      # Progress Tracking tests
│   ├── repositories/        # Data access layer (to be implemented)
│   └── utils/               # Utility functions (to be implemented)
├── test/
│   ├── fixtures/            # Test data builders and predefined data
│   │   └── fixtures.go      # Test fixtures for all domain models
│   ├── helpers/             # Test helper functions
│   │   └── test_helpers.go  # Common test utilities
│   └── reporters/           # Custom Ginkgo reporters
│       └── gherkin_reporter.go  # Given-When-Then style reporter
├── suite_test.go            # Main test suite configuration
├── Makefile                 # Test automation commands
└── README_TESTS.md          # This file
```

## Test Philosophy

We follow **Test-Driven Development (TDD)** principles:

1. **Red**: Write a failing test
2. **Green**: Write minimal code to make it pass
3. **Refactor**: Clean up while keeping tests green

### Given-When-Then Style

All tests use BDD-style naming with Ginkgo:

```go
Context("Given a task on the critical path", func() {
    Context("And the task is due within 3 days", func() {
        When("priority is calculated", func() {
            It("should have priority score above 90", func() {
                // Test code
            })
        })
    })
})
```

## Running Tests

### Prerequisites

```bash
# Install Ginkgo CLI
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Install dependencies
go mod download
```

### Basic Commands

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test file
go test -v ./internal/services/priority_service_test.go

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Using Makefile

```bash
# Run all tests
make test

# Run with coverage report
make test-coverage

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Run tests in watch mode
make test-watch

# Run all CI checks
make ci
```

## Test Modules

### 1. Priority Engine (`priority_service_test.go`)

Tests the universal priority calculation algorithm:
- Critical path task prioritization
- Deadline urgency scoring
- Dependency impact calculation
- Project priority weighting
- Task ranking

**Key Test Cases:**
- Tasks on critical path with upcoming deadlines
- Overdue task handling
- Multiple dependency impact
- Priority score bounds (0-100)

### 2. Health Scoring (`health_service_test.go`)

Tests project and portfolio health calculations:
- Portfolio health aggregation
- Schedule health (progress vs timeline)
- Completion health
- Dependency health
- Resource health
- Critical path health

**Key Test Cases:**
- Behind-schedule project detection
- Resource overallocation detection
- Health trend analysis
- Status categorization (healthy/caution/at-risk/critical)

### 3. Nudge Detection (`nudge_service_test.go`)

Tests the 7 nudge type detection algorithms:
- Overload detection
- Delay risk detection
- Skill gap detection
- Unassigned task detection
- Blocked task detection
- Resource conflict detection
- Dependency block detection

**Key Test Cases:**
- Severity classification
- Criticality scoring
- Nudge deduplication
- Rate limiting

### 4. Assignment Engine (`assignment_service_test.go`)

Tests smart task assignment:
- Skill match scoring (40 points max)
- Availability scoring (30 points max)
- Workload balance scoring (20 points max)
- Past performance scoring (10 points max)
- Context switch penalty

**Key Test Cases:**
- Best candidate selection
- Overallocation warnings
- Skill gap identification

### 5. Dependency Management (`dependency_service_test.go`)

Tests task dependency management:
- Circular dependency detection
- Critical path calculation (CPM)
- Topological sorting
- Float/slack calculation
- Impact analysis

**Key Test Cases:**
- Self-dependency rejection
- Cycle detection in diamond patterns
- Zero-float critical path identification

### 6. Scenario Processing (`scenario_service_test.go`)

Tests what-if scenario simulation:
- Employee leave impact
- Scope change impact
- Reallocation impact
- Priority shift impact
- Timeline recalculation
- Cost analysis

**Key Test Cases:**
- Cascade effect detection
- Before/after comparison
- Recommendation generation

### 7. Progress Tracking (`progress_service_test.go`)

Tests project and task progress tracking:
- Completion percentage calculation
- Hierarchical progress roll-up
- Progress variance tracking
- Velocity calculation
- Milestone tracking

**Key Test Cases:**
- Weighted completion by hours
- Status-to-progress mapping
- Behind/ahead schedule detection

## Test Data

### Fixtures

The `test/fixtures` package provides:

- **Builders**: Fluent API for creating test objects
  ```go
  user := fixtures.NewUser().
      WithID("user-123").
      WithName("John Doe").
      WithRole(models.RoleAdmin).
      Build()
  ```

- **Predefined Data**:
  - `CreateTestTeam()` - 8 realistic team members
  - `CreateTestSkills()` - 12 skills across categories
  - `CreateTestProjects()` - 3 active projects
  - `CreateTestTasks()` - 12+ tasks in various states
  - `CreateTestWorkloadData()` - Allocation data

### Custom Reporter

The `GherkinReporter` outputs tests in Given-When-Then format:

```
🧪 Test Suite: Priority Engine Service
============================================================

✅ Scenario: should have priority score above 90
   Given a task on the critical path
   And the task is due within 3 days
   When priority is calculated
   Then it should have priority score above 90 (12ms)

✅ Scenario: should return error for nil task
   Given invalid task data
   When priority is calculated
   Then it should return error for nil task (5ms)

📊 Test Results:
   ✅ Passed:  25
   ❌ Failed:  0
   ⏭️  Skipped: 0
   📈 Total:   25
```

## Writing New Tests

### 1. Create Test File

```bash
# Using Ginkgo CLI
cd internal/services
ginkgo generate my_service

# Or manually
touch my_service_test.go
```

### 2. Test File Template

```go
package services_test

import (
    "testing"
    
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    
    "github.com/xephyr-ai/xephyr-backend/test/fixtures"
    "github.com/xephyr-ai/xephyr-backend/test/helpers"
)

func TestMyService(t *testing.T) {
    helpers.RunSuite(t, "My Service")
}

var _ = Describe("My Feature", func() {
    
    Context("Given some initial state", func() {
        var data SomeType
        
        BeforeEach(func() {
            data = fixtures.NewSomething().Build()
        })
        
        When("an action occurs", func() {
            It("should produce expected result", func() {
                result := DoSomething(data)
                Expect(result).To(Equal(expected))
            })
        })
    })
})
```

### 3. Use Table-Driven Tests for Multiple Cases

```go
DescribeTable("boundary condition tests",
    func(input int, expected int) {
        result := CalculateSomething(input)
        Expect(result).To(Equal(expected))
    },
    Entry("minimum value", 0, 0),
    Entry("normal value", 50, 100),
    Entry("maximum value", 100, 200),
)
```

## Continuous Integration

The `make ci` command runs:
1. `go fmt` - Code formatting
2. `go vet` - Static analysis
3. `golangci-lint run` - Linting
4. `go test -coverprofile` - Tests with coverage

## Coverage Goals

| Module | Target Coverage |
|--------|----------------|
| Priority Engine | 90% |
| Health Scoring | 90% |
| Nudge Detection | 85% |
| Assignment Engine | 85% |
| Dependency Management | 90% |
| Scenario Processing | 80% |
| Progress Tracking | 85% |

## Next Steps

1. **Implement Service Layer**: Write actual implementations to make tests pass
2. **Add Repository Tests**: Create database integration tests
3. **Add API Tests**: Test HTTP handlers and middleware
4. **Add E2E Tests**: Full workflow testing

## Resources

- [Ginkgo Documentation](https://onsi.github.io/ginkgo/)
- [Gomega Matchers](https://onsi.github.io/gomega/)
- [Go Testing Best Practices](https://github.com/golang/go/wiki/TestComments)
