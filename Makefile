# Makefile for SageMaker Cost Calculator

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=mohua

# Directories
SRC_DIR=.
INTERNAL_DIR=internal

# Build flags
LDFLAGS=-s -w

# Default target
all: test build

# Install UPX if not present
install-upx:
	@if ! command -v upx >/dev/null 2>&1; then \
		if [ "$$(uname)" = "Darwin" ]; then \
			brew install upx; \
		elif [ "$$(uname)" = "Linux" ]; then \
			sudo apt-get update && sudo apt-get install -y upx-ucl; \
		else \
			echo "Please install UPX manually on Windows using: choco install upx"; \
			exit 1; \
		fi \
	fi

# Build the binary with optimizations
build: install-upx
	@echo "Building optimized binary..."
	CGO_ENABLED=0 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) $(SRC_DIR)
	@echo "Compressing with UPX..."
	upx --best --lzma $(BINARY_NAME)
	@echo "Build complete"

# Run only unit tests (default)
test:
	$(GOTEST) -v -tags=!integration ./...

# Run integration tests (requires AWS credentials)
test-integ:
	$(GOTEST) -v -tags=integration ./...

# Run all tests (both unit and integration)
test-all:
	$(GOTEST) -v ./... -tags=integration

# Run tests with coverage (unit tests only)
cover:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html

# Tidy dependencies
deps:
	$(GOMOD) tidy
	$(GOMOD) verify

# Lint the code
lint:
	golangci-lint run

# Install development tools
install-dev-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) $(SRC_DIR)
	./$(BINARY_NAME)

.PHONY: all build test cover clean deps lint install-dev-tools run
