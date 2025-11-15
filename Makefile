.PHONY: help build clean test examples run-examples baml fmt lint install

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=pres
EXAMPLE_DIR=examples
BUILD_DIR=build
GO=go
BAML=baml-cli

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

help: ## Show this help message
	@echo "$(COLOR_BOLD)Available targets:$(COLOR_RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-15s$(COLOR_RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(COLOR_BOLD)Examples:$(COLOR_RESET)"
	@echo "  make build          # Build the main binary"
	@echo "  make examples       # Build all examples"
	@echo "  make run-examples   # Run the question components demo"
	@echo "  make clean          # Remove build artifacts"
	@echo ""

build: ## Build the main pres binary
	@echo "$(COLOR_BLUE)Building $(BINARY_NAME)...$(COLOR_RESET)"
	@$(GO) build -o $(BINARY_NAME) .
	@echo "$(COLOR_GREEN)✓ Built $(BINARY_NAME)$(COLOR_RESET)"

examples: ## Build all example programs
	@echo "$(COLOR_BLUE)Building examples...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/input_demo $(EXAMPLE_DIR)/input_components.go
	@echo "$(COLOR_GREEN)✓ Built $(BUILD_DIR)/input_demo$(COLOR_RESET)"

run-examples: examples ## Build and run the input components example
	@echo "$(COLOR_BLUE)Running input components demo...$(COLOR_RESET)"
	@echo ""
	@./$(BUILD_DIR)/input_demo

install: build ## Install the binary to $GOPATH/bin
	@echo "$(COLOR_BLUE)Installing $(BINARY_NAME)...$(COLOR_RESET)"
	@$(GO) install .
	@echo "$(COLOR_GREEN)✓ Installed $(BINARY_NAME) to $(shell go env GOPATH)/bin$(COLOR_RESET)"

test: ## Run tests
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	@$(GO) test -v ./...
	@echo "$(COLOR_GREEN)✓ Tests passed$(COLOR_RESET)"

test-coverage: ## Run tests with coverage
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	@$(GO) test -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✓ Coverage report generated: coverage.html$(COLOR_RESET)"

baml: ## Generate BAML client code
	@echo "$(COLOR_BLUE)Generating BAML client...$(COLOR_RESET)"
	@$(BAML) generate
	@echo "$(COLOR_GREEN)✓ BAML client generated$(COLOR_RESET)"

fmt: ## Format Go code
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	@$(GO) fmt ./...
	@gofmt -s -w .
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

lint: ## Run linters
	@echo "$(COLOR_BLUE)Running linters...$(COLOR_RESET)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
		echo "$(COLOR_GREEN)✓ Linting complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ golangci-lint not found, skipping...$(COLOR_RESET)"; \
		echo "  Install with: brew install golangci-lint"; \
	fi

tidy: ## Tidy go modules
	@echo "$(COLOR_BLUE)Tidying modules...$(COLOR_RESET)"
	@$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Modules tidied$(COLOR_RESET)"

clean: ## Remove build artifacts
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -rf presentations/*.json presentations/*.html
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

clean-all: clean ## Remove all generated files including BAML client
	@echo "$(COLOR_BLUE)Removing BAML client...$(COLOR_RESET)"
	@rm -rf baml_client
	@echo "$(COLOR_GREEN)✓ Deep clean complete$(COLOR_RESET)"

rebuild: clean build ## Clean and rebuild

all: clean baml build examples test ## Clean, generate BAML, build everything, and test

dev: ## Development mode - build and run with a test presentation
	@echo "$(COLOR_BLUE)Development mode...$(COLOR_RESET)"
	@$(MAKE) build
	@echo ""
	@echo "$(COLOR_YELLOW)Try these commands:$(COLOR_RESET)"
	@echo "  ./$(BINARY_NAME) create \"Test Presentation\""
	@echo "  ./$(BINARY_NAME) generate --path presentations/test-presentation.json"
	@echo "  ./$(BINARY_NAME) update --path presentations/test-presentation.json \"Add a conclusion slide\""
	@echo ""

check: fmt lint test ## Run formatters, linters, and tests

watch: ## Watch for changes and rebuild (requires entr)
	@if command -v entr > /dev/null; then \
		echo "$(COLOR_BLUE)Watching for changes...$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Press Ctrl+C to stop$(COLOR_RESET)"; \
		find . -name '*.go' -not -path './baml_client/*' | entr -c make build; \
	else \
		echo "$(COLOR_YELLOW)⚠ entr not found$(COLOR_RESET)"; \
		echo "  Install with: brew install entr"; \
	fi

version: ## Show version information
	@echo "$(COLOR_BOLD)pres - Presentation Generator$(COLOR_RESET)"
	@echo "Go version: $(shell $(GO) version)"
	@echo "BAML version: $(shell $(BAML) --version 2>/dev/null || echo 'not found')"

deps: ## Show dependencies
	@echo "$(COLOR_BLUE)Project dependencies:$(COLOR_RESET)"
	@$(GO) list -m all

update-deps: ## Update dependencies
	@echo "$(COLOR_BLUE)Updating dependencies...$(COLOR_RESET)"
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies updated$(COLOR_RESET)"

.PHONY: build-all
build-all: build examples ## Build main binary and all examples
	@echo "$(COLOR_GREEN)✓ All binaries built$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Built artifacts:$(COLOR_RESET)"
	@ls -lh $(BINARY_NAME) 2>/dev/null || true
	@ls -lh $(BUILD_DIR)/* 2>/dev/null || true
