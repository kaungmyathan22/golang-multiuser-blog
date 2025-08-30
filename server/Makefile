# Makefile for Go Blog Application

# Variables
APP_NAME = blog-server
BINARY_NAME = blog-server
MAIN_FILE = cmd/server/main.go

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Docker parameters
DOCKER_CMD = docker
DOCKER_COMPOSE = docker-compose

# Default target
all: build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)
	./$(BINARY_NAME)

# Install dependencies
deps:
	$(GOMOD) download

# Run unit tests
test-unit:
	$(GOTEST) -v ./internal/... -short

# Run integration tests
test-integration:
	$(GOTEST) -v ./internal/... -run Integration

# Run end-to-end tests
test-e2e:
	$(GOTEST) -v ./internal/e2e/... -tags=e2e

# Run all tests
test-all:
	$(GOTEST) -v ./internal/... -tags=e2e

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./internal/...
	$(GOCMD) tool cover -html=coverage.out

# Docker build
docker-build:
	$(DOCKER_CMD) build -t $(APP_NAME) .

# Run with docker-compose
docker-run:
	$(DOCKER_COMPOSE) up --build

# Stop docker-compose
docker-stop:
	$(DOCKER_COMPOSE) down

# Run docker-compose in detached mode
docker-start:
	$(DOCKER_COMPOSE) up --build -d

# View docker-compose logs
docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Help
help:
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@echo "  all              Build the application (default)"
	@echo "  build            Build the application"
	@echo "  run              Build and run the application"
	@echo "  deps             Install dependencies"
	@echo "  test-unit        Run unit tests"
	@echo "  test-integration Run integration tests"
	@echo "  test-e2e         Run end-to-end tests"
	@echo "  test-all         Run all tests"
	@echo "  clean            Clean build files"
	@echo "  fmt              Format code"
	@echo "  vet              Vet code"
	@echo "  test-coverage    Run tests with coverage"
	@echo "  docker-build     Build Docker image"
	@echo "  docker-run       Run with docker-compose"
	@echo "  docker-stop      Stop docker-compose"
	@echo "  docker-start     Run docker-compose in detached mode"
	@echo "  docker-logs      View docker-compose logs"
	@echo "  help             Show this help message"

.PHONY: all build run deps test-unit test-integration test-e2e test-all clean fmt vet test-coverage docker-build docker-run docker-stop docker-start docker-logs help