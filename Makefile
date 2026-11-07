.PHONY: all build run cli test fmt vet lint check install ci docker docker-build docker-up docker-down up down start stop redis-up redis-down dev prod logs docker-logs restart status clean clean-all help

GO ?= go
DOCKER ?= docker
COMPOSE ?= $(DOCKER) compose
OUT_DIR := bin
APP_NAME := einoagent
CLI_NAME := einoagentcli
VERSION ?= dev
IMAGE_NAME ?= eino_assistant
IMAGE_TAG ?= $(VERSION)
ENV_FILE ?= .env
COMPOSE_FILE ?= docker-compose.yml
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
GOLANGCI_LINT ?= golangci-lint

all: build

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build Go binaries into $(OUT_DIR)"
	@echo "  run           Run the main agent app (loads $(ENV_FILE) if present)"
	@echo "  cli           Run the CLI app"
	@echo "  test          Run go test ./..."
	@echo "  fmt           Format Go source files"
	@echo "  vet           Run go vet ./..."
	@echo "  lint          Run golangci-lint if installed"
	@echo "  check         Run fmt + vet + test"
	@echo "  install       Download Go modules and install lint tools"
	@echo "  ci            Run install + check + build"
	@echo "  docker        Build Docker image (alias for docker-build)"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-up     Start Docker Compose services"
	@echo "  docker-down   Stop Docker Compose services"
	@echo "  up            Alias for docker-up"
	@echo "  down          Alias for docker-down"
	@echo "  start         Alias for docker-up"
	@echo "  stop          Alias for docker-down"
	@echo "  redis-up      Start Redis service"
	@echo "  redis-down    Stop Redis service"
	@echo "  dev           Start Redis + build + run app"
	@echo "  prod          Build image and show production-style run hint"
	@echo "  logs          Show Docker Compose logs"
	@echo "  docker-logs   Alias for logs"
	@echo "  restart       Restart Docker Compose services"
	@echo "  status        Show Docker Compose status"
	@echo "  clean         Remove build artifacts"
	@echo "  clean-all     Remove build artifacts and Docker cache (optional)"
	@echo "  release       Build and package the app for release"

load-env:
	@if [ -f "$(ENV_FILE)" ]; then \
		set -a; . "$(ENV_FILE)"; set +a; \
	fi

install:
	$(GO) mod download
	@if command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		echo "golangci-lint already installed"; \
	else \
		$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi

build: $(OUT_DIR)/$(APP_NAME) $(OUT_DIR)/$(CLI_NAME)

$(OUT_DIR):
	@mkdir -p $(OUT_DIR)

$(OUT_DIR)/$(APP_NAME): | $(OUT_DIR)
	$(GO) build $(LDFLAGS) -o $(OUT_DIR)/$(APP_NAME) ./cmd/einoagent

$(OUT_DIR)/$(CLI_NAME): | $(OUT_DIR)
	$(GO) build $(LDFLAGS) -o $(OUT_DIR)/$(CLI_NAME) ./cmd/einoagentcli

run: load-env $(OUT_DIR)/$(APP_NAME)
	./$(OUT_DIR)/$(APP_NAME)

cli: load-env $(OUT_DIR)/$(CLI_NAME)
	./$(OUT_DIR)/$(CLI_NAME)

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

lint:
	@if command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		$(GOLANGCI_LINT) run; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		false; \
	fi

check: fmt vet test

ci: install check build

release: build docker-build
	@echo "Release assets generated in $(OUT_DIR) and image tag: $(IMAGE_NAME):$(IMAGE_TAG)"

docker: docker-build

docker-build:
	$(DOCKER) build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-up:
	$(COMPOSE) -f $(COMPOSE_FILE) up -d

docker-down:
	$(COMPOSE) -f $(COMPOSE_FILE) down

up: docker-up

down: docker-down

start: docker-up

stop: docker-down

redis-up: docker-up
redis-down: docker-down

logs:
	$(COMPOSE) -f $(COMPOSE_FILE) logs -f

docker-logs: logs

restart:
	$(COMPOSE) -f $(COMPOSE_FILE) restart

status:
	$(COMPOSE) -f $(COMPOSE_FILE) ps

dev: docker-up build run

prod: docker-build
	@echo "Production image ready: $(IMAGE_NAME):$(IMAGE_TAG)"
	@echo "Run with: $(DOCKER) run --rm --env-file $(ENV_FILE) -p 8080:8080 $(IMAGE_NAME):$(IMAGE_TAG)"

clean:
	rm -rf $(OUT_DIR)

clean-all: clean
	$(DOCKER) system df >/dev/null 2>&1 || true
