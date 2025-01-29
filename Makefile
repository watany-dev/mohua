.PHONY: all build clean install-upx

# Get the version information
VERSION := $(shell git describe --tags --always)
COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%d')
LDFLAGS := -s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${BUILD_DATE}

# Default target
all: build

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
	CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o sagemaker-monitor
	@echo "Compressing with UPX..."
	upx --best --lzma sagemaker-monitor
	@echo "Build complete: sagemaker-monitor"

# Build for all platforms
build-all: install-upx
	@echo "Building for all platforms..."
	# Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o sagemaker-monitor-linux
	upx --best --lzma sagemaker-monitor-linux
	# macOS
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o sagemaker-monitor-darwin
	upx --best --lzma sagemaker-monitor-darwin
	# Windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o sagemaker-monitor-windows.exe
	upx --best --lzma sagemaker-monitor-windows.exe
	@echo "Build complete for all platforms"

# Clean build artifacts
clean:
	rm -f sagemaker-monitor sagemaker-monitor-linux sagemaker-monitor-darwin sagemaker-monitor-windows.exe

# Show help
help:
	@echo "Available targets:"
	@echo "  make          - Build optimized binary for current platform"
	@echo "  make build    - Same as above"
	@echo "  make build-all- Build optimized binaries for all platforms"
	@echo "  make clean    - Remove built binaries"
	@echo "  make help     - Show this help"
