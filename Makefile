.PHONY: help build run-server run-cli test clean install deps

# Default target
help:
	@echo "IamFeel - Development Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make deps          Install Go dependencies"
	@echo "  make install       Install binaries to PATH"
	@echo ""
	@echo "Development:"
	@echo "  make build         Build all binaries"
	@echo "  make run-server    Run the web server"
	@echo "  make run-cli       Run CLI (use ARGS='command' for specific commands)"
	@echo ""
	@echo "Testing:"
	@echo "  make test          Run all tests"
	@echo "  make test-verbose  Run tests with verbose output"
	@echo "  make coverage      Generate test coverage report"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean         Remove build artifacts"
	@echo "  make fmt           Format code"
	@echo "  make lint          Run linter"
	@echo ""
	@echo "Examples:"
	@echo "  make run-cli ARGS='onboard'    Run onboarding"
	@echo "  make run-cli ARGS='plan'       Generate plan"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Build binaries
build:
	@echo "Building CLI..."
	go build -o bin/iamfeel cmd/cli/*.go
	@echo "Building server..."
	go build -o bin/server cmd/server/main.go
	@echo "✓ Build complete! Binaries in ./bin/"

# Install to system PATH
install: build
	@echo "Installing to system..."
	go install ./cmd/server
	go install ./cmd/cli

# Run web server
run-server:
	@echo "Starting web server..."
	go run cmd/server/main.go

# Run CLI (use ARGS to pass commands)
run-cli:
	@echo "Running CLI..."
	go run cmd/cli/main.go $(ARGS)

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	go test -v ./...

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Create necessary directories
setup-dirs:
	@echo "Creating necessary directories..."
	mkdir -p data bin
	@echo "Directories created!"
