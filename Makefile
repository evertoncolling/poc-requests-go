# Makefile for poc-requests-go

.PHONY: help build test test-verbose test-race test-coverage clean lint fmt vet deps check run install-tools proto-gen

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Build parameters
BINARY_NAME=poc-requests-go
MAIN_PATH=./
BUILD_DIR=./bin

# Coverage parameters
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Default target
help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the application
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Test targets
test: ## Run unit tests (excludes integration tests)
	@echo "Running unit tests..."
	$(GOTEST) -v -short ./...

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	$(GOTEST) -v -short ./...

test-integration: ## Run integration tests (requires .env file with credentials)
	@echo "Running integration tests..."
	@if [ -f .env ]; then \
		echo "Loading environment from .env file..."; \
		set -a && source .env && set +a && $(GOTEST) -v -run TestIntegration ./pkg/api/; \
	else \
		echo "No .env file found, checking environment variables..."; \
		if [ -z "$$CLIENT_ID" ] || [ -z "$$CLIENT_SECRET" ] || [ -z "$$TENANT_ID" ] || [ -z "$$CDF_CLUSTER" ] || [ -z "$$CDF_PROJECT" ]; then \
			echo "❌ Integration tests require credentials. Create .env file or set environment variables."; \
			exit 1; \
		fi; \
		$(GOTEST) -v -run TestIntegration ./pkg/api/; \
	fi

test-integration-ci: ## Run integration tests for CI (uses environment variables directly)
	@echo "Running integration tests for CI..."
	$(GOTEST) -v -run TestIntegration ./pkg/api/

test-all: ## Run all tests including integration tests (requires credentials)
	@echo "Running all tests..."
	@if [ -f .env ]; then \
		echo "Loading environment from .env file..."; \
		set -a && source .env && set +a && $(GOTEST) -v ./...; \
	else \
		echo "No .env file found, checking environment variables..."; \
		if [ -z "$$CLIENT_ID" ] || [ -z "$$CLIENT_SECRET" ] || [ -z "$$TENANT_ID" ] || [ -z "$$CDF_CLUSTER" ] || [ -z "$$CDF_PROJECT" ]; then \
			echo "❌ Integration tests require credentials. Use 'make test-unit' for unit tests only."; \
			exit 1; \
		fi; \
		$(GOTEST) -v ./...; \
	fi

test-short: ## Run tests with short flag
	@echo "Running short tests..."
	$(GOTEST) -short -v ./...

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	$(GOTEST) -race -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

test-coverage-func: ## Show test coverage by function
	@echo "Running tests with coverage by function..."
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

# Code quality targets
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

lint-fix: ## Run linter with autofix
	@echo "Running linter with autofix..."
	golangci-lint run --fix

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) -s -w .

fmt-check: ## Check if code is formatted
	@echo "Checking code formatting..."
	@if [ -n "$$($(GOFMT) -l .)" ]; then \
		echo "Code is not formatted. Run 'make fmt'"; \
		$(GOFMT) -l .; \
		exit 1; \
	fi

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

check: fmt-check vet lint test ## Run all code quality checks

# Dependency targets
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download

deps-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	$(GOMOD) verify

deps-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Development targets
run: ## Run the application
	@echo "Running application..."
	$(GOCMD) run $(MAIN_PATH)

install-tools: ## Install development tools
	@echo "Installing development tools..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) -u google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GOGET) -u github.com/google/go-licenses@latest

proto-gen: ## Generate protobuf files
	@echo "Generating protobuf files..."
	@if [ ! -d "pkg/proto" ]; then \
		echo "proto directory not found"; \
		exit 1; \
	fi
	@find pkg/proto -name "*.proto" -exec protoc --go_out=. --go_opt=paths=source_relative {} \;

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

clean-cache: ## Clean module cache
	@echo "Cleaning module cache..."
	$(GOCMD) clean -modcache

# Docker targets (if needed in future)
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm -it $(BINARY_NAME)

# Documentation targets
docs: ## Generate documentation
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation server started at http://localhost:6060"

# Security targets
security: ## Run security scan
	@echo "Running security scan..."
	gosec ./...

# Benchmarks
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# All-in-one targets
all: clean deps fmt vet lint test build ## Run all tasks

ci: fmt-check vet lint test ## Run CI pipeline locally