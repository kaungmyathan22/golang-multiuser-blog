#!/bin/bash

# Test script for running tests with Docker

set -e

echo "Starting test databases..."
docker-compose -f docker-compose.test.yml up -d

# Wait for databases to be ready
echo "Waiting for databases to be ready..."
sleep 10

# Run tests
echo "Running unit tests..."
go test -v ./internal/... -short

echo "Running integration tests..."
go test -v ./internal/... -run Integration

echo "Running end-to-end tests..."
go test -v ./internal/e2e/... -tags=e2e

# Stop test databases
echo "Stopping test databases..."
docker-compose -f docker-compose.test.yml down

echo "All tests completed!"