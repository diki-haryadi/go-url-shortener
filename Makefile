# Variables
BINARY_NAME=shortener-api
DOCKER_COMPOSE=docker-compose
GOFLAGS=-ldflags="-w -s"

# Default shell
SHELL=/bin/bash

.PHONY: all build clean test coverage deps docker-build docker-up docker-down help

# Default target
all: help

## Build:
build: ## Build the application
	@echo "Building..."
	go build ${GOFLAGS} -o ${BINARY_NAME} cmd/main.go

clean: ## Clean up built binary
	@echo "Cleaning..."
	rm -f ${BINARY_NAME}
	go clean

## Test:
test: ## Run tests
	go test -v ./...

coverage: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## Dependencies:
deps: ## Download dependencies
	go mod download
	go mod tidy

## Docker:
docker-build: ## Build docker image
	${DOCKER_COMPOSE} build

docker-up: ## Start docker containers
	${DOCKER_COMPOSE} up -d

docker-down: ## Stop docker containers
	${DOCKER_COMPOSE} down

docker-logs: ## View docker logs
	${DOCKER_COMPOSE} logs -f

## Development:
dev: ## Run application in development mode
	go run cmd/main.go

lint: ## Run linter
	golangci-lint run

## Database:
redis-cli: ## Connect to Redis CLI
	docker exec -it shortener-api_redis_1 redis-cli

## Deployment:
deploy-prod: ## Deploy to production
	@echo "Deploying to production..."
	# Add your deployment commands here

## Help:
help: ## Show this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Additional useful targets
docker-rebuild: docker-down docker-build docker-up ## Rebuild and restart containers

# Migrations if needed
migrate-up: ## Run database migrations
	@echo "Running migrations..."
	# Add your migration commands here

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	# Add your migration commands here

# Testing specific targets
test-integration: ## Run integration tests
	go test -tags=integration ./...

test-coverage-html: ## Generate HTML coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html