# Makefile for SageMaker Cost Calculator

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=sagemaker-cost-calculator

# Directories
SRC_DIR=.
INTERNAL_DIR=internal

# Build flags
LDFLAGS=-s -w

# Default target
all: test build

# Build the application
build:
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) $(SRC_DIR)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
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
