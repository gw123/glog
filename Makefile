.PHONY: test test-verbose test-cover test-race bench clean fmt vet lint help install-tools

# Go parameters
GOROOT?=~/.gvm/gos/go1.21
GOPATH?=$(shell go env GOPATH 2>/dev/null || echo ~/go)
GOCMD=$(shell which go || echo $(GOROOT)/bin/go)
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Test parameters
TEST_FLAGS=-v -count=1
COVER_FLAGS=-coverprofile=coverage.out -covermode=atomic
RACE_FLAGS=-race
BENCH_FLAGS=-bench=. -benchmem -benchtime=1s -run=^$$

# Directories
PKG_LIST=$(shell go list ./... | grep -v /vendor/)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run all tests
	@echo "Running tests..."
	$(GOTEST) ./... $(TEST_FLAGS)

test-short: ## Run tests without verbose output
	@echo "Running tests (short)..."
	$(GOTEST) ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	$(GOTEST) ./... -v

test-cover: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) ./... $(COVER_FLAGS)
	@echo "Coverage report generated: coverage.out"
	@$(GOCMD) tool cover -func=coverage.out | tail -1

test-cover-html: test-cover ## Generate HTML coverage report
	@echo "Generating HTML coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	$(GOTEST) ./... $(RACE_FLAGS)

test-all: test-race test-cover ## Run all test suites

bench: ## Run benchmark tests
	@echo "Running benchmark tests..."
	$(GOTEST) ./... $(BENCH_FLAGS)

bench-cpu: ## Run benchmark tests with CPU profiling
	@echo "Running benchmark tests with CPU profiling..."
	$(GOTEST) ./... $(BENCH_FLAGS) -cpuprofile=cpu.out
	@echo "CPU profile generated: cpu.out"

bench-mem: ## Run benchmark tests with memory profiling
	@echo "Running benchmark tests with memory profiling..."
	$(GOTEST) ./... $(BENCH_FLAGS) -memprofile=mem.out
	@echo "Memory profile generated: mem.out"

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

lint: ## Run linters (requires golangci-lint)
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: make install-tools" && exit 1)
	golangci-lint run ./...

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	$(GOMOD) tidy

download: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download

verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	$(GOMOD) verify

clean: ## Clean build cache and test files
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f coverage.out coverage.html
	rm -f cpu.out mem.out
	rm -f *.log test*.log
	rm -rf test_logs/

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@which golangci-lint > /dev/null || \
		(echo "Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin)
	@echo "Development tools installed!"

check: fmt vet ## Run format and vet checks

ci: clean tidy check test-race test-cover ## Run CI pipeline locally

demo: ## Run demo application
	@echo "Running demo..."
	$(GOCMD) run ./demo/basic/main.go

build-demo: ## Build demo application
	@echo "Building demo..."
	$(GOBUILD) -o bin/basic-demo ./demo/basic/main.go
	$(GOBUILD) -o bin/webapp-demo ./demo/webapp/main.go
	$(GOBUILD) -o bin/performance-demo ./demo/performance/main.go

# Quality checks
quality: fmt vet lint test-cover ## Run all quality checks

# Quick test
quick: ## Quick test (no race detector, no coverage)
	@echo "Running quick tests..."
	$(GOTEST) ./... -short

# Watch mode (requires entr)
watch: ## Watch for changes and run tests (requires entr)
	@which entr > /dev/null || (echo "entr not installed. Install it first." && exit 1)
	@echo "Watching for changes..."
	@find . -name '*.go' | entr -c make test-short

# Statistics
stats: ## Show code statistics
	@echo "=== Code Statistics ==="
	@echo "Total lines of code:"
	@find . -name '*.go' -not -path "./vendor/*" | xargs wc -l | tail -1
	@echo "\nNumber of Go files:"
	@find . -name '*.go' -not -path "./vendor/*" | wc -l
	@echo "\nNumber of test files:"
	@find . -name '*_test.go' -not -path "./vendor/*" | wc -l
	@echo "\nTest coverage:"
	@make test-cover 2>&1 | grep "coverage:" || echo "Run 'make test-cover' first"

# Default target
all: check test ## Run checks and tests

.DEFAULT_GOAL := help
