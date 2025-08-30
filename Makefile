# Golang Multi-User Blog Makefile (Monolith)
# Variables
APP_NAME = golang-multiuser-blog
BUILD_DIR = bin
SERVER_DIR = server
CGO_ENABLED = 0

# Colors for output
RESET = \033[0m
BOLD = \033[1m
RED = \033[31m
GREEN = \033[32m
YELLOW = \033[33m
BLUE = \033[34m
MAGENTA = \033[35m
CYAN = \033[36m

# Default target
.PHONY: help
help: ## Show this help message
	@echo "$(BOLD)$(APP_NAME) - Available Commands:$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""

# Setup and installation
.PHONY: install
install: ## Install project dependencies and tools
	@echo "$(YELLOW)Installing Go dependencies...$(RESET)"
	@cd $(SERVER_DIR) && go mod download
	@echo "$(YELLOW)Installing air for hot reload...$(RESET)"
	@go install github.com/cosmtrek/air@latest
	@echo "$(YELLOW)Installing golangci-lint for code quality...$(RESET)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)Setup complete!$(RESET)"

# Clean
.PHONY: clean
clean: ## Clean build artifacts and temporary files
	@rm -rf $(BUILD_DIR)
	@rm -rf $(SERVER_DIR)/tmp
	@echo "$(GREEN)Cleaned build artifacts$(RESET)"

# Development with hot reload
.PHONY: dev
dev: ## Run server with hot reload
	@echo "$(CYAN)Starting server with hot reload...$(RESET)"
	@cd $(SERVER_DIR) && air -c .air.toml

# Build commands
.PHONY: build
build: ## Build server
	@echo "$(YELLOW)Building server...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@cd $(SERVER_DIR) && CGO_ENABLED=$(CGO_ENABLED) go build -o ../$(BUILD_DIR)/$(APP_NAME) ./cmd/server
	@echo "$(GREEN)Server built successfully$(RESET)"

# Run
.PHONY: run
run: build ## Build and run server
	@echo "$(CYAN)Running server...$(RESET)"
	@./$(BUILD_DIR)/$(APP_NAME)

# Test commands
.PHONY: test
test: ## Run all tests
	@echo "$(YELLOW)Running tests...$(RESET)"
	@cd $(SERVER_DIR) && go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "$(YELLOW)Running tests with coverage...$(RESET)"
	@cd $(SERVER_DIR) && go test -v -coverprofile=coverage.out ./...
	@cd $(SERVER_DIR) && go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(RESET)"

.PHONY: test-watch
test-watch: ## Run tests with file watching
	@echo "$(CYAN)Starting test watcher...$(RESET)"
	@cd $(SERVER_DIR) && air -c .air.test.toml

# Code quality
.PHONY: lint
lint: ## Run linter
	@echo "$(YELLOW)Running linter...$(RESET)"
	@cd $(SERVER_DIR) && golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	@echo "$(YELLOW)Formatting code...$(RESET)"
	@cd $(SERVER_DIR) && go fmt ./...
	@echo "$(GREEN)Code formatted$(RESET)"

# Database
.PHONY: migrate-up
migrate-up: ## Run database migrations
	@echo "$(YELLOW)Running database migrations...$(RESET)"
	@cd $(SERVER_DIR) && ./scripts/migrate.sh up

.PHONY: migrate-down
migrate-down: ## Rollback database migrations
	@echo "$(YELLOW)Rolling back database migrations...$(RESET)"
	@cd $(SERVER_DIR) && ./scripts/migrate.sh down

# Docker
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(RESET)"
	@docker build -f docker/Dockerfile -t $(APP_NAME) .
	@echo "$(GREEN)Docker image built successfully$(RESET)"

.PHONY: docker-run
docker-run: docker-build ## Build and run with Docker
	@echo "$(CYAN)Running with Docker...$(RESET)"
	@docker run -p 8080:8080 $(APP_NAME)

# Utilities
.PHONY: mod-tidy
mod-tidy: ## Tidy Go modules
	@echo "$(YELLOW)Tidying Go modules...$(RESET)"
	@cd $(SERVER_DIR) && go mod tidy
	@echo "$(GREEN)Go modules tidied$(RESET)"

.PHONY: gen
gen: ## Run code generation
	@echo "$(YELLOW)Running code generation...$(RESET)"
	@cd $(SERVER_DIR) && go generate ./...
	@echo "$(GREEN)Code generation complete$(RESET)"