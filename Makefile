.PHONY: build test clean install lint fmt vet

# Variables
BINARY_NAME=githubber
VERSION=2.0.0
BUILD_DIR=build
MAIN_PATH=./cmd

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

# Install to system PATH
install-system: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Uninstall from system PATH
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Lint the code
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format the code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet the code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Run all checks
check: fmt vet lint test

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@go mod download
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run the application
run: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Watch for changes and rebuild
watch:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

# Generate documentation
docs:
	@echo "Generating documentation..."
	@go doc -all > docs/API.md

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .

# Docker run
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run --rm -it $(BINARY_NAME):$(VERSION)

# Release
release: clean build-all
	@echo "Creating release $(VERSION)..."
	@mkdir -p release
	@cp $(BUILD_DIR)/* release/
	@tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-amd64
	@tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-amd64
	@tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-arm64
	@zip release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip -j $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-all      - Build for all platforms"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install to GOPATH/bin"
	@echo "  install-system - Install to /usr/local/bin"
	@echo "  uninstall      - Remove from system PATH"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  check          - Run all checks"
	@echo "  dev-setup      - Setup development environment"
	@echo "  run            - Build and run"
	@echo "  watch          - Watch for changes and rebuild"
	@echo "  docs           - Generate documentation"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  release        - Create release packages"
	@echo "  help           - Show this help"