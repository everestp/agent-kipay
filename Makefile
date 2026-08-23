# ==============================================================================
# VARIABLES
# ==============================================================================
APP_NAME = bheri-backend
BUILD_DIR = bin
MAIN_PATH = cmd/main.go

# Load environment variables from a .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# ==============================================================================
# COMMANDS
# ==============================================================================

.PHONY: all run build clean fmt lint test

all: clean build

# Run the application locally
run:
	@echo "Starting $(APP_NAME)..."
	go run $(MAIN_PATH)

# Build the binary into the bin/ directory
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# Format Go code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet for static analysis
lint:
	@echo "Running linter..."
	go vet ./...

# Run tests
test:
	@echo "Running tests..."
	go test -v -race ./...
