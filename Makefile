.PHONY: all build test lint clean run-api run-cli fmt vet tidy

BINARY_NAME=carrot-agent
API_BINARY=carrot-agent-api
CLI_BINARY=carrot-agent-cli
GO=go
GOFLAGS=-ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

VERSION ?= 0.1.0
BUILD_TIME ?= $$(date -u '+%Y-%m-%d %H:%M:%S')

all: tidy fmt vet test build

build:
	$(GO) build $(GOFLAGS) -o bin/$(API_BINARY) ./cmd/api
	$(GO) build $(GOFLAGS) -o bin/$(CLI_BINARY) ./cmd/cli

build-api:
	$(GO) build $(GOFLAGS) -o bin/$(API_BINARY) ./cmd/api

build-cli:
	$(GO) build $(GOFLAGS) -o bin/$(CLI_BINARY) ./cmd/cli

test:
	$(GO) test -v -race ./...

test-unit:
	$(GO) test -v -short ./...

test-cover:
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

run-api:
	$(GO) run ./cmd/api

run-cli:
	$(GO) run ./cmd/cli

docker-build:
	docker build -t carrotagent/carrot-agent:latest .

docker-run:
	docker run -it carrotagent/carrot-agent:latest

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

install-deps:
	$(GO) get ./...

check: tidy fmt vet lint test

help:
	@echo "Available targets:"
	@echo "  all          - Run tidy, fmt, vet, test, and build"
	@echo "  build        - Build API and CLI binaries"
	@echo "  build-api    - Build API binary only"
	@echo "  build-cli    - Build CLI binary only"
	@echo "  test         - Run all tests with race detection"
	@echo "  test-unit    - Run unit tests only (short mode)"
	@echo "  test-cover   - Run tests with coverage report"
	@echo "  lint         - Run golangci-lint"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  tidy         - Tidy modules"
	@echo "  clean        - Remove build artifacts"
	@echo "  run-api      - Run API server"
	@echo "  run-cli      - Run CLI client"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  check        - Run all checks (tidy, fmt, vet, lint, test)"