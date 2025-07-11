# Auto-Sort Downloads - Makefile

BINARY_NAME=auto-sort
BUILD_DIR=build
MAIN_FILE=cmd/main.go

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)
	
	@echo "Multi-platform build complete!"

# Build for Windows only
build-windows:
	@echo "Building for Windows..."
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@echo "Windows build complete: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

# Build for macOS only
build-macos:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	@echo "macOS build complete!"

# Build for Linux only
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)
	@echo "Linux build complete!"

# Build with optimizations (smaller binary size)
build-optimized:
	@echo "Building optimized binary..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-optimized $(MAIN_FILE)
	@echo "Optimized build complete: $(BUILD_DIR)/$(BINARY_NAME)-optimized"

# Build with debug information
build-debug:
	@echo "Building debug binary..."
	@mkdir -p $(BUILD_DIR)
	go build -gcflags="all=-N -l" -o $(BUILD_DIR)/$(BINARY_NAME)-debug $(MAIN_FILE)
	@echo "Debug build complete: $(BUILD_DIR)/$(BINARY_NAME)-debug"

# Run the application
run:
	go run $(MAIN_FILE)

# Run in dry-run mode
dry-run:
	go run $(MAIN_FILE) -dry-run

# Generate default configuration
config:
	go run $(MAIN_FILE) -generate-config

# Install dependencies
deps:
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f auto-sort.log
	@echo "Clean complete!"

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Install the binary to system PATH
install: build
	@echo "Installing $(BINARY_NAME) to system PATH..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete! You can now run '$(BINARY_NAME)' from anywhere."

# Uninstall the binary from system PATH
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstall complete!"

# Show help
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make build-all     - Build for multiple platforms"
	@echo "  make build-windows - Build for Windows only"
	@echo "  make build-macos   - Build for macOS only"
	@echo "  make build-linux   - Build for Linux only"
	@echo "  make build-optimized - Build optimized binary (smaller size)"
	@echo "  make build-debug   - Build with debug information"
	@echo "  make run           - Run the application"
	@echo "  make dry-run     - Run in dry-run mode"
	@echo "  make config      - Generate default configuration"
	@echo "  make deps        - Install dependencies"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make install     - Install to system PATH"
	@echo "  make uninstall   - Remove from system PATH"
	@echo "  make help        - Show this help"

.PHONY: all build build-all build-windows build-macos build-linux build-optimized build-debug run dry-run config deps clean test fmt lint install uninstall help