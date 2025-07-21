# GitHubber - Makefile
# Author: Ritankar Saha <ritankar.saha786@gmail.com>

.PHONY: build install uninstall test clean run help dev

APP_NAME := githubber
BUILD_DIR := ./build
BINARY_PATH := $(BUILD_DIR)/$(APP_NAME)
INSTALL_PATH := /usr/local/bin/$(APP_NAME)

# Default target
all: build

# Build the application
build:
	@echo "🔨 Building GitHubber..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.version=$$(git describe --tags --always --dirty)" -o $(BINARY_PATH) ./cmd/main.go
	@echo "✅ Build complete: $(BINARY_PATH)"

# Install the application globally
install: build
	@echo "📦 Installing GitHubber..."
	@sudo cp $(BINARY_PATH) $(INSTALL_PATH)
	@sudo chmod +x $(INSTALL_PATH)
	@echo "✅ GitHubber installed to $(INSTALL_PATH)"
	@echo "🚀 You can now run 'githubber' from anywhere!"

# Uninstall the application
uninstall:
	@echo "🗑️  Uninstalling GitHubber..."
	@sudo rm -f $(INSTALL_PATH)
	@echo "✅ GitHubber uninstalled"

# Run the application (for development)
run:
	@echo "🚀 Running GitHubber..."
	@go run ./cmd/main.go

# Run in development mode with hot reload
dev:
	@echo "🔧 Starting development mode..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Lint the code
lint:
	@echo "🔍 Linting code..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run

# Format the code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@goimports -w .

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

# Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	@go mod download
	@go mod tidy

# Cross-compile for multiple platforms
build-all:
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./cmd/main.go
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/main.go
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./cmd/main.go
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/main.go
	@echo "✅ Cross-compilation complete"

# Create release archive
release: build-all
	@echo "📦 Creating release archives..."
	@cd $(BUILD_DIR) && tar -czf $(APP_NAME)-linux-amd64.tar.gz $(APP_NAME)-linux-amd64
	@cd $(BUILD_DIR) && tar -czf $(APP_NAME)-darwin-amd64.tar.gz $(APP_NAME)-darwin-amd64
	@cd $(BUILD_DIR) && tar -czf $(APP_NAME)-darwin-arm64.tar.gz $(APP_NAME)-darwin-arm64
	@cd $(BUILD_DIR) && zip $(APP_NAME)-windows-amd64.zip $(APP_NAME)-windows-amd64.exe
	@echo "✅ Release archives created in $(BUILD_DIR)"

# Check if we're in a git repository for version info
check-git:
	@git status > /dev/null 2>&1 || (echo "❌ Not a git repository" && exit 1)

# Show help
help:
	@echo "GitHubber - Build System"
	@echo "========================"
	@echo ""
	@echo "Available targets:"
	@echo "  build        Build the application"
	@echo "  install      Install the application globally"
	@echo "  uninstall    Remove the installed application"
	@echo "  run          Run the application (development)"
	@echo "  dev          Run in development mode with hot reload"
	@echo "  test         Run tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  lint         Lint the code"
	@echo "  fmt          Format the code"
	@echo "  clean        Clean build artifacts"
	@echo "  deps         Download and tidy dependencies"
	@echo "  build-all    Cross-compile for multiple platforms"
	@echo "  release      Create release archives"
	@echo "  help         Show this help message"
	@echo ""
	@echo "Example usage:"
	@echo "  make build      # Build the application"
	@echo "  make install    # Build and install globally"
	@echo "  make test       # Run tests"
	@echo "  make clean      # Clean build artifacts"