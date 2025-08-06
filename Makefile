.PHONY: build run test clean docker-build docker-up docker-down docker-logs help

# Variables
APP_NAME := transaction-service
DOCKER_COMPOSE := docker-compose
GO_FILES := $(shell find . -name "*.go" -type f)

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build commands
build: ## Build the application binary
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) .

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	@go run main.go

# Testing commands
test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@$(DOCKER_COMPOSE) build

docker-up: ## Start services with Docker Compose
	@echo "Starting services..."
	@$(DOCKER_COMPOSE) up -d

docker-down: ## Stop services
	@echo "Stopping services..."
	@$(DOCKER_COMPOSE) down

docker-down-volumes: ## Stop services and remove volumes
	@echo "Stopping services and removing volumes..."
	@$(DOCKER_COMPOSE) down -v

docker-logs: ## Show Docker logs
	@$(DOCKER_COMPOSE) logs -f

docker-restart: ## Restart services
	@echo "Restarting services..."
	@$(DOCKER_COMPOSE) restart

# Development commands
dev: docker-up ## Start development environment
	@echo "Development environment started!"
	@echo "API available at: http://localhost:8080"
	@echo "Use 'make docker-logs' to view logs"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@$(DOCKER_COMPOSE) down -v --remove-orphans
	@docker system prune -f

# Database commands
db-reset: ## Reset database (remove volumes and restart)
	@echo "Resetting database..."
	@$(DOCKER_COMPOSE) down -v
	@$(DOCKER_COMPOSE) up -d postgres
	@sleep 5
	@$(DOCKER_COMPOSE) up -d app

# Utility commands
check: ## Check code formatting and run basic validations
	@echo "Checking code formatting..."
	@gofmt -l $(GO_FILES)
	@echo "Running go vet..."
	@go vet ./...
	@echo "Running go mod tidy..."
	@go mod tidy

format: ## Format Go code
	@echo "Formatting code..."
	@gofmt -w $(GO_FILES)

# Quick test commands
test-balance: ## Test balance endpoint for user 1
	@echo "Testing balance endpoint..."
	@curl -s http://localhost:8080/user/1/balance | json_pp || echo "Service might not be running"

test-transaction: ## Test transaction endpoint with sample data
	@echo "Testing transaction endpoint..."
	@curl -s -X POST http://localhost:8080/user/1/transaction \
		-H "Source-Type: game" \
		-H "Content-Type: application/json" \
		-d '{"state":"win","amount":"10.50","transactionId":"test-'$(shell date +%s)'"}' | json_pp || echo "Service might not be running"

load-test: ## Run simple load test (requires hey: go install github.com/rakyll/hey@latest)
	@echo "Running load test..."
	@hey -n 100 -c 10 -m GET http://localhost:8080/user/1/balance