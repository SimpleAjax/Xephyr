# Xephyr Backend Makefile

.PHONY: test test-unit test-integration test-coverage test-report clean help

# Default test command
.DEFAULT_GOAL := test

## help: Show this help message
help:
	@echo "Available commands:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //'

## test: Run all tests with Ginkgo
 test:
	go test -v ./... -count=1

## test-unit: Run unit tests only
 test-unit:
	go test -v ./internal/services/... -count=1 -label=unit

## test-integration: Run integration tests
 test-integration:
	go test -v ./internal/... -count=1 -label=integration

## test-coverage: Run tests with coverage report
 test-coverage:
	go test -v ./... -coverprofile=coverage.out -count=1
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## test-report: Run tests with custom Gherkin report
 test-report:
	go test -v ./... -count=1 -json > test-report.json
	@echo "Test report generated: test-report.json"

## test-watch: Run tests in watch mode (requires ginkgo CLI)
 test-watch:
	ginkgo watch -v ./...

## ginkgo-install: Install Ginkgo CLI
 ginkgo-install:
	go install github.com/onsi/ginkgo/v2/ginkgo@latest

## ginkgo-bootstrap: Bootstrap Ginkgo test suite
 ginkgo-bootstrap:
	ginkgo bootstrap

## ginkgo-generate: Generate new test file (usage: make ginkgo-generate PKG=services)
 ginkgo-generate:
	@if [ -z "$(PKG)" ]; then \
		echo "Usage: make ginkgo-generate PKG=packagename"; \
		exit 1; \
	fi
	cd internal/$(PKG) && ginkgo generate

## lint: Run Go linter
 lint:
	golangci-lint run ./...

## fmt: Format Go code
 fmt:
	go fmt ./...

## vet: Run go vet
 vet:
	go vet ./...

## deps: Download dependencies
 deps:
	go mod download
	go mod tidy

## deps-update: Update all dependencies
 deps-update:
	go get -u ./...
	go mod tidy

## clean: Clean test artifacts
 clean:
	rm -f coverage.out coverage.html test-report.json
	rm -rf test-results/

## db-test-up: Start test database (Docker)
 db-test-up:
	docker run -d \
		--name xephyr-test-db \
		-e POSTGRES_DB=xephyr_test \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-p 5433:5432 \
		postgres:15-alpine
	@echo "Waiting for database to be ready..."
	@sleep 3

## db-test-down: Stop test database
 db-test-down:
	docker stop xephyr-test-db || true
	docker rm xephyr-test-db || true

## db-test-reset: Reset test database
 db-test-reset: db-test-down db-test-up

## ci: Run all CI checks
 ci: fmt vet lint test-coverage
