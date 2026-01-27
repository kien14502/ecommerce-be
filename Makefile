APP_NAME := ecommerce
BINARY_NAME := ecommerce-be
GO_VERSION := 1.25.6
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := latest
PORT := 8080
ROOT_GO := cmd/server/main.go
GO_CMD := go

DOCKER_PRODUCT := docker.prod.compose.yml
DOCKER_DEV := docker.compose.yml
DOCKER_TEST := docker.test.compose.yml

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	mkdir -p bin
	$(GOBUILD) -o bin/$(BINARY_NAME) -v $(ROOT_GO)

# Run targets
run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	$(GO_CMD) run $(ROOT_GO)

run-dev: ## Run the application in development mode
	@echo "Running $(BINARY_NAME) in development mode..."
	$(GO_CMD) run $(ROOT_GO) -e development

run-prod: ## Run the application in production mode
	@echo "Running $(BINARY_NAME) in production mode..."
	$(GO_CMD) run $(ROOT_GO) -e production

# 	Docker targets
docker-dev: ## Run the application in a Docker container for development
	@echo "Running $(DOCKER_IMAGE) in Docker container for development..."
	docker-compose -f developments/docker/$(DOCKER_DEV) up --build

docker-prod: ## Run the application in a Docker container for production
	@echo "Running $(DOCKER_IMAGE) in Docker container for production..."
	docker-compose -f developments/docker/$(DOCKER_PRODUCT) up --build

docker-test: ## Run the application in a Docker container for testing
	@echo "Running $(DOCKER_IMAGE) in Docker container for testing..."
	docker-compose -f developments/docker/$(DOCKER_TEST) up --build