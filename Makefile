# FiscaFlow Makefile
# Provides common development tasks and shortcuts

.PHONY: help test test-unit test-integration test-coverage test-benchmark test-race lint security clean build run docker-build docker-run docker-stop docker-clean docker-dev docker-prod docker-logs docker-logs-all docker-restart docker-health

# Default target
help:
	@echo "FiscaFlow Development Commands"
	@echo "=============================="
	@echo ""
	@echo "Testing:"
	@echo "  test              - Run all tests"
	@echo "  test-unit         - Run unit tests only"
	@echo "  test-integration  - Run integration tests only"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  test-benchmark    - Run benchmarks"
	@echo "  test-race         - Run tests with race detection"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint              - Run linting"
	@echo "  security          - Run security scan"
	@echo "  format            - Format code with gofmt"
	@echo "  tidy              - Tidy Go modules"
	@echo ""
	@echo "Build & Run:"
	@echo "  build             - Build the application"
	@echo "  run               - Run the application"
	@echo "  clean             - Clean build artifacts"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build      - Build Docker image"
	@echo "  docker-run        - Run with Docker Compose"
	@echo "  docker-stop       - Stop Docker Compose services"
	@echo "  docker-clean      - Clean Docker resources"
	@echo "  docker-dev        - Start development environment with hot reloading"
	@echo "  docker-prod       - Start production environment"
	@echo "  docker-logs       - View FiscaFlow logs"
	@echo "  docker-logs-all   - View all service logs"
	@echo "  docker-restart    - Restart FiscaFlow service"
	@echo "  docker-health     - Check service health"
	@echo ""
	@echo "Development:"
	@echo "  deps              - Install dependencies"
	@echo "  migrate           - Run database migrations"
	@echo "  seed              - Seed database with test data"

# Testing targets
test:
	@echo "Running all tests..."
	@./scripts/test.sh all

test-unit:
	@echo "Running unit tests..."
	@./scripts/test.sh unit

test-integration:
	@echo "Running integration tests..."
	@./scripts/test.sh integration

test-coverage:
	@echo "Running tests with coverage..."
	@./scripts/test.sh coverage

test-benchmark:
	@echo "Running benchmarks..."
	@./scripts/test.sh benchmark

test-race:
	@echo "Running tests with race detection..."
	@./scripts/test.sh race

# Code quality targets
lint:
	@echo "Running linting..."
	@./scripts/test.sh lint

security:
	@echo "Running security scan..."
	@./scripts/test.sh security

format:
	@echo "Formatting code..."
	@gofmt -s -w .

tidy:
	@echo "Tidying Go modules..."
	@go mod tidy

# Build targets
build:
	@echo "Building application..."
	@go build -o bin/server cmd/server/main.go

run:
	@echo "Running application..."
	@go run cmd/server/main.go

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage_*.out coverage_*.html
	@go clean

# Docker targets
docker-build:
	@echo "[INFO] Building Docker image..."
	docker build -t fiscaflow:latest .

docker-run:
	@echo "[INFO] Starting FiscaFlow with Docker Compose..."
	docker-compose up --build

docker-run-d:
	@echo "[INFO] Starting FiscaFlow with Docker Compose (detached)..."
	docker-compose up --build -d

docker-stop:
	@echo "[INFO] Stopping Docker Compose services..."
	docker-compose down

docker-clean:
	@echo "[INFO] Cleaning Docker resources..."
	docker-compose down -v --remove-orphans
	docker system prune -f

docker-dev:
	@echo "[INFO] Starting development environment with hot reloading..."
	docker-compose -f docker-compose.dev.yml up --build -d

docker-prod:
	@echo "[INFO] Starting production environment..."
	docker-compose -f docker-compose.prod.yml up --build

docker-logs:
	@echo "[INFO] Viewing FiscaFlow logs..."
	docker-compose logs -f fiscaflow

docker-logs-all:
	@echo "[INFO] Viewing all service logs..."
	docker-compose logs -f

docker-restart:
	@echo "[INFO] Restarting FiscaFlow service..."
	docker-compose restart fiscaflow

docker-health:
	@echo "[INFO] Checking service health..."
	docker-compose ps

# Development targets
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

migrate:
	@echo "Running database migrations..."
	@# TODO: Implement migration command

seed:
	@echo "Seeding database..."
	@# TODO: Implement seed command

# CI/CD targets
ci-test:
	@echo "Running CI tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

ci-build:
	@echo "Building for CI..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go

# Development shortcuts
dev: deps tidy lint test-unit run

quick-test: test-unit test-integration

full-test: lint security test test-race test-benchmark test-coverage 