# Makefile cho User Service Tests

.PHONY: test test-unit test-integration test-all test-coverage test-benchmark clean help

# Default target
help:
	@echo "Available commands:"
	@echo "  test            - Run all tests"
	@echo "  test-unit       - Run unit tests only"
	@echo "  test-integration- Run integration tests only"
	@echo "  test-all        - Run all tests with verbose output"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-benchmark  - Run benchmark tests"
	@echo "  clean           - Clean test artifacts"
	@echo "  setup-test-db   - Setup test database"

# Run all tests
test:
	go test ./... -v

# Run unit tests only (mock-based)
test-unit:
	go test ./... -v -tags=unit

# Run integration tests only (database-based)
test-integration:
	go test ./... -v -tags=integration

# Run all tests with verbose output
test-all:
	go test ./... -v -count=1

# Run tests with coverage report
test-coverage:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmark tests
test-benchmark:
	go test ./... -v -bench=. -benchmem

# Clean test artifacts
clean:
	rm -f coverage.out coverage.html
	go clean -testcache

# Setup test database
setup-test-db:
	@echo "Setting up test database..."
	@echo "Make sure PostgreSQL is running and create database 'user_service_test'"
	@echo "Example:"
	@echo "  createdb -U postgres user_service_test"
	@echo "  psql -U postgres -d user_service_test -c 'CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";'"

# Run tests with race detection
test-race:
	go test ./... -v -race

# Run tests with specific timeout
test-timeout:
	go test ./... -v -timeout=30s

# Run tests and generate JUnit XML report (for CI/CD)
test-junit:
	go install github.com/jstemmer/go-junit-report@latest
	go test ./... -v 2>&1 | go-junit-report > test-report.xml

# Run tests with specific package
test-package:
	@read -p "Enter package name: " package; \
	go test ./$$package -v

# Run tests with specific function
test-function:
	@read -p "Enter test function name: " function; \
	go test ./... -v -run $$function

# Run tests with specific tags
test-tags:
	@read -p "Enter tags (e.g., unit,integration): " tags; \
	go test ./... -v -tags=$$tags

# Install test dependencies
install-test-deps:
	go get github.com/stretchr/testify/assert
	go get github.com/stretchr/testify/require
	go get github.com/stretchr/testify/suite
	go get github.com/stretchr/testify/mock

# Show test help
test-help:
	@echo "Test Environment Variables:"
	@echo "  TEST_DATABASE_URL - Database connection string for tests"
	@echo "  TEST_LOG_LEVEL   - Log level for tests (default: silent)"
	@echo ""
	@echo "Example:"
	@echo "  TEST_DATABASE_URL='host=localhost user=postgres password=postgres dbname=user_service_test port=5432 sslmode=disable' make test-integration"

# Run tests in parallel
test-parallel:
	go test ./... -v -parallel=4

# Run tests with memory profiling
test-memprofile:
	go test ./... -v -memprofile=mem.prof
	go tool pprof -top mem.prof

# Run tests with CPU profiling
test-cpuprofile:
	go test ./... -v -cpuprofile=cpu.prof
	go tool pprof -top cpu.prof 