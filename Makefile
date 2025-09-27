.PHONY: help run build test clean sqlc-generate migrate-up migrate-down docker-up docker-down

# Variables
APP_NAME=go-ddd-clean-starter
BINARY_DIR=bin
BINARY_NAME=$(BINARY_DIR)/api
MAIN_PATH=cmd/api/main.go

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

help: ## Show this help message
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

run: ## Run the application
	@echo "$(GREEN)Starting application...$(NC)"
	go run $(MAIN_PATH)

build: ## Build the application binary
	@echo "$(GREEN)Building application...$(NC)"
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Binary created at $(BINARY_NAME)$(NC)"

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Run tests with coverage report
	@echo "$(GREEN)Generating coverage report...$(NC)"
	go tool cover -html=coverage.txt -o coverage.html
	@echo "$(GREEN)Coverage report generated at coverage.html$(NC)"

clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning...$(NC)"
	rm -rf $(BINARY_DIR)
	rm -f coverage.txt coverage.html
	@echo "$(GREEN)Clean complete$(NC)"

sqlc-generate: ## Generate SQLC code
	@echo "$(GREEN)Generating SQLC code...$(NC)"
	sqlc generate
	@echo "$(GREEN)SQLC generation complete$(NC)"

sqlc-generate-go: ## Generate SQLC code using go run (no sqlc installation needed)
	@echo "$(GREEN)Generating SQLC code using go run...$(NC)"
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate
	@echo "$(GREEN)SQLC generation complete$(NC)"

migrate-up: ## Run database migrations up
	@echo "$(GREEN)Running migrations up...$(NC)"
	@echo "$(YELLOW)Note: Install golang-migrate first: https://github.com/golang-migrate/migrate$(NC)"
	migrate -path infrastructure/database/migrations -database "postgresql://postgres:postgres@localhost:5432/go_ddd_starter?sslmode=disable" up

migrate-down: ## Run database migrations down
	@echo "$(GREEN)Running migrations down...$(NC)"
	migrate -path infrastructure/database/migrations -database "postgresql://postgres:postgres@localhost:5432/go_ddd_starter?sslmode=disable" down

migrate-create: ## Create a new migration file (usage: make migrate-create name=create_users_table)
	@echo "$(GREEN)Creating migration: $(name)$(NC)"
	migrate create -ext sql -dir infrastructure/database/migrations -seq $(name)

docker-up: ## Start Docker containers (PostgreSQL)
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Containers started$(NC)"

docker-down: ## Stop Docker containers
	@echo "$(GREEN)Stopping Docker containers...$(NC)"
	docker-compose down
	@echo "$(GREEN)Containers stopped$(NC)"

docker-logs: ## Show Docker container logs
	docker-compose logs -f

deps: ## Download dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)Dependencies downloaded$(NC)"

fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)Code formatted$(NC)"

lint: ## Run linter
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run ./...

dev: ## Run in development mode with hot reload (requires air)
	@echo "$(GREEN)Starting development server with hot reload...$(NC)"
	air

install-tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)Tools installed$(NC)"
