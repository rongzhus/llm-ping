BINARY_NAME=llm-ping
VERSION=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%d\ %H:%M:%S)
GO_VERSION=$(shell go version)

LDFLAGS=-ldflags "-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'"

.PHONY: all build build-all clean help test install

all: build

build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) main.go
	@echo "Build complete: $(BINARY_NAME)"

build-all: build-linux build-windows build-macos
	@echo "All builds complete!"

build-linux:
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 main.go
	@echo "Built: $(BINARY_NAME)-linux-amd64"

build-linux-arm64:
	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 main.go
	@echo "Built: $(BINARY_NAME)-linux-arm64"

build-windows:
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe main.go
	@echo "Built: $(BINARY_NAME)-windows-amd64.exe"

build-macos:
	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 main.go
	@echo "Built: $(BINARY_NAME)-darwin-amd64"

build-macos-arm64:
	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 main.go
	@echo "Built: $(BINARY_NAME)-darwin-arm64"

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	@echo "Clean complete"

install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) main.go
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

help:
	@echo "Available targets:"
	@echo "  all            - Build for current platform (default)"
	@echo "  build          - Build for current platform"
	@echo "  build-all      - Build for all platforms (Linux, Windows, macOS)"
	@echo "  build-linux    - Build for Linux (amd64)"
	@echo "  build-linux-arm64 - Build for Linux (arm64)"
	@echo "  build-windows  - Build for Windows (amd64)"
	@echo "  build-macos    - Build for macOS (amd64)"
	@echo "  build-macos-arm64 - Build for macOS (arm64)"
	@echo "  test           - Run tests"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install to GOPATH/bin"
	@echo "  help           - Show this help message"
