# Makefile for Gin Simple REST API

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=gin-app
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out coverage.html

# Run the application
run:
	$(GOCMD) run ./cmd/server

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Test the application
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -cover ./...

# Run tests with coverage report
test-coverage-html:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run specific test
test-func:
	@read -p "Enter test function name: " func; \
	$(GOTEST) -v -run $$func ./...

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Start server in background
start:
	$(GOCMD) run ./cmd/server &

# Stop all go processes (be careful with this)
stop:
	pkill -f "go run ./cmd/server" || true

# Install additional tools
install-tools:
	$(GOGET) -u github.com/swaggo/swag/cmd/swag

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  clean           - Clean build files"
	@echo "  run             - Run the application"
	@echo "  deps            - Download dependencies"
	@echo "  test            - Run unit tests (Go tests)"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  test-coverage-html - Generate HTML coverage report"
	@echo "  test-func       - Run specific test function"
	@echo "  start           - Start server in background"
	@echo "  stop            - Stop background server"
	@echo "  fmt             - Format code"
	@echo "  lint            - Lint code"
	@echo "  help            - Show this help message"

.PHONY: build clean run deps test test-coverage test-coverage-html test-func start stop install-tools fmt lint help
