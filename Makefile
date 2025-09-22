# Go CLI Makefile

# Variables
GO = /Users/milos-dev/go/go1.25.1/bin/go
BINARY_NAME = go-cli
BUILD_DIR = bin

# ANSI color codes
GREEN = \033[0;32m
RED = \033[0;31m
NC = \033[0m # No Color

.PHONY: build run test test-integration clean help

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/main.go
	@echo "$(GREEN)Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run tests with colored output
test:
	@echo "Running unit tests..."
	@$(GO) test -v ./... | sed 's/--- PASS:/--- \x1b[0;32mPASS:\x1b[0m/g; s/--- FAIL:/--- \x1b[0;31mFAIL:\x1b[0m/g; s/^PASS$$/\x1b[0;32mPASS\x1b[0m/g; s/^FAIL$$/\x1b[0;31mFAIL\x1b[0m/g'

# Run integration tests with real repositories
test-integration: build
	@echo "Running integration tests..."
	@./scripts/integration-test.sh

# Setup integration test environment and run sync command
integration-demo: build
	@echo "$(GREEN)Setting up integration test environment...$(NC)"
	@mkdir -p /tmp/go-cli-integration
	@echo "Cloning test repositories..."
	@cd /tmp/go-cli-integration && \
		git clone --quiet https://github.com/octocat/Hello-World.git hello-world 2>/dev/null || true && \
		git clone --quiet https://github.com/github/gitignore.git gitignore 2>/dev/null || true && \
		git clone --quiet https://github.com/torvalds/linux.git --depth 1 linux 2>/dev/null || true
	@echo "$(GREEN)Running sync command on test repositories...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME) sync -d /tmp/go-cli-integration || echo "$(YELLOW)Some repositories may have failed (this is expected for demo)$(NC)"
	@echo "$(GREEN)Integration demo completed!$(NC)"
	@echo "Cleaning up temporary directory..."
	@rm -rf /tmp/go-cli-integration
	@echo "$(GREEN)Cleanup completed$(NC)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)Clean completed$(NC)"

# Show help
help:
	@echo "Available commands:"
	@echo "  build            - Build the application"
	@echo "  run              - Build and run the application"
	@echo "  test             - Run unit tests"
	@echo "  test-integration - Run integration tests with script"
	@echo "  integration-demo - Demo integration test with real repos"
	@echo "  clean            - Clean build artifacts"
	@echo "  help             - Show this help message"