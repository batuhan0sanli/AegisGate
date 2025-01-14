# Build variables
BINARY_NAME=aegisgate
DOCKER_REGISTRY=
DOCKER_IMAGE=aegisgate
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH=$(shell git rev-parse --short HEAD)

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Docker related variables
DOCKERFILE=build/Docker/Dockerfile
DOCKER_COMPOSE_FILE=deployments/docker-compose.yml

# Make settings
.DEFAULT_GOAL := help
.PHONY: all build clean test coverage docker-build docker-push run help

all: clean build test ## Build and run tests

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags="-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(COMMIT_HASH)'" \
		-o $(GOBIN)/$(BINARY_NAME) ./cmd/aegisgate

clean: ## Clean build directory
	@echo "Cleaning..."
	@rm -rf $(GOBIN)
	@go clean
	@rm -f coverage.out

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run

docker-build: ## Build docker image
	@echo "Building docker image..."
	docker build -t $(DOCKER_REGISTRY)$(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_REGISTRY)$(DOCKER_IMAGE):latest \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg COMMIT_HASH=$(COMMIT_HASH) \
		-f $(DOCKERFILE) .

docker-push: ## Push docker image to registry
	@echo "Pushing docker image..."
	docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE):$(VERSION)

run: build ## Run the application
	@echo "Running $(BINARY_NAME)..."
	@$(GOBIN)/$(BINARY_NAME)

docker-compose-up: ## Start the application using docker-compose
	@echo "Starting services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) -p $(BINARY_NAME) up -d

docker-compose-down: ## Stop the application using docker-compose
	@echo "Stopping services..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) -p $(BINARY_NAME) down

docker-compose-logs: ## View docker-compose logs
	@echo "Showing logs..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) -p $(BINARY_NAME) logs -f

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' 