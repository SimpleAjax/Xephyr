# Xephyr E2E Tests

End-to-end tests for the Xephyr API, following the TDD Guide and API Design specifications.

## Overview

These E2E tests validate complete user flows through the API, ensuring all components work together correctly. The tests use Ginkgo v2 with Gomega for BDD-style testing.

## Test Structure

```
tests/e2e/
├── suite_test.go                 # Main test suite setup
├── helpers/
│   ├── api_client.go            # HTTP client for API requests
│   └── response_parser.go       # Response parsing utilities
├── flows/
│   ├── task_assignment_flow_test.go     # Task assignment scenarios
│   ├── scenario_simulation_flow_test.go # Scenario simulation scenarios
│   ├── dependency_management_flow_test.go # Dependency management scenarios
│   ├── health_monitoring_flow_test.go   # Health monitoring scenarios
│   ├── nudge_handling_flow_test.go      # Nudge handling scenarios
│   └── critical_paths_smoke_test.go     # Critical path smoke tests
└── README.md
```

## Running Tests

### Run all E2E tests
```bash
cd backend
go test ./tests/e2e/... -v
```

### Run specific test file
```bash
go test ./tests/e2e/flows -v -run "Task Assignment Flow"
```

### Run with custom reporter (Gherkin format)
```bash
go test ./tests/e2e/... -v -ginkgo.json-report=report.json
```

### Run against external server
```bash
E2E_BASE_URL=https://api.xephyr.io E2E_API_TOKEN=your-token go test ./tests/e2e/... -v
```

## Test Configuration

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `E2E_BASE_URL` | Base URL for API (empty = start test server) | "" |
| `E2E_API_TOKEN` | API authentication token | "test-token" |
| `E2E_ORG_ID` | Organization ID for requests | "org-test-123" |
| `E2E_SKIP_CLEANUP` | Skip test data cleanup | "false" |

## Test Categories

### Flow Tests
Complete user journeys:
- **Task Assignment Flow**: Creating tasks, getting suggestions, assigning
- **Scenario Simulation Flow**: Creating scenarios, running simulations, applying changes
- **Dependency Management Flow**: Creating dependencies, validating, checking critical path
- **Health Monitoring Flow**: Portfolio health, project health, trends
- **Nudge Handling Flow**: Listing nudges, taking actions, statistics

### Smoke Tests
Quick validation of critical paths:
- Priority API endpoints
- Health API endpoints
- Assignment API endpoints
- Dependency API endpoints
- Nudge API endpoints
- Complete PM decision-making workflow

## Writing New Tests

Follow the Given-When-Then pattern:

```go
Describe("Given [initial context]", func() {
    Context("When [action]", func() {
        It("should [expected outcome]", func() {
            // Arrange
            req := dto.SomeRequest{...}
            
            // Act
            resp, err := client.Post("/endpoint", req)
            
            // Assert
            Expect(err).ToNot(HaveOccurred())
            Expect(resp.StatusCode).To(Equal(http.StatusOK))
        })
    })
})
```

## Test Data

Tests use fixtures from `test/fixtures/fixtures.go`:
- Test team (Sarah, Mike, Alex, Emma, James, etc.)
- Test projects (E-Commerce, Fitness App, SaaS Dashboard)
- Test tasks with various states
- Test skills and user skill mappings

## Guidelines

1. **Independent Tests**: Each test should be independent
2. **Cleanup**: Use `AfterEach` to clean up created resources
3. **Timeouts**: Tests have 60-second timeout by default
4. **Retry Logic**: Use `DoWithRetry` for flaky operations
5. **Descriptive Names**: Use clear Given-When-Then descriptions
