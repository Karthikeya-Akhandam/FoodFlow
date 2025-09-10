# FoodFlow Makefile

.PHONY: help build run test clean docker-build docker-run migrate-up migrate-down seed sqlc-generate keys

# Default target
help: ## Show this help message
	@echo "FoodFlow - Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development commands
build: ## Build the application
	@echo "Building FoodFlow API..."
	go build -o bin/foodflow ./cmd/api

run: ## Run the application locally
	@echo "Running FoodFlow API..."
	go run ./cmd/api

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Database commands (auto-migration runs on app startup)
# Migrations run automatically when the application starts

seed: ## Seed the database with sample data
	@echo "Seeding database..."
	@echo "Note: Auto-migration runs automatically when you start the app with 'make run'"

# Code generation (removed sqlc - using direct SQL queries)

# JWT Keys
keys: ## Generate JWT RSA keys
	@echo "Generating JWT RSA keys..."
	mkdir -p keys
	openssl genrsa -out keys/app.rsa 2048
	openssl rsa -in keys/app.rsa -pubout > keys/app.rsa.pub
	chmod 600 keys/app.rsa
	chmod 644 keys/app.rsa.pub

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t foodflow:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env foodflow:latest

docker-compose-up: ## Start all services with docker-compose
	@echo "Starting all services..."
	docker-compose up -d

docker-compose-down: ## Stop all services
	@echo "Stopping all services..."
	docker-compose down

docker-compose-logs: ## View logs from all services
	@echo "Viewing logs..."
	docker-compose logs -f

docker-compose-restart: ## Restart all services
	@echo "Restarting all services..."
	docker-compose restart

# Development setup
setup: ## Setup development environment
	@echo "Setting up development environment..."
	go run scripts/setup.go

# Production commands
prod-build: ## Build production Docker image
	@echo "Building production Docker image..."
	docker build -t foodflow:prod -f Dockerfile.prod .

# Utility commands
fmt: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy

# Database utilities
db-reset: ## Reset database (WARNING: This will delete all data)
	@echo "Database reset complete! Auto-migration will run on next app startup."

db-shell: ## Connect to database shell
	@echo "Connecting to database shell..."
	psql postgres://foodflow:password@localhost:5432/foodflow

# Monitoring
logs: ## View application logs
	@echo "Viewing application logs..."
	docker-compose logs -f api

health: ## Check application health
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || echo "Application is not healthy"
