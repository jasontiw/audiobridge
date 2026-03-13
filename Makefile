# Makefile for AudioBridge

# Build variables
BINARY_NAME=audiobridge
GO=go
GOOS=$(shell $(GO) env GOOS)
GOARCH=$(shell $(GO) env GOARCH)

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME)

# Build for all platforms
.PHONY: build-all
build-all:
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BINARY_NAME)-darwin-arm64
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-linux-amd64
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY_NAME).exe

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Run tests
.PHONY: test
test:
	$(GO) test -v -race ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html

# Run go mod tidy
.PHONY: mod-tidy
mod-tidy:
	$(GO) mod tidy

# Install dependencies
.PHONY: deps
deps:
	$(GO) mod download

# Show help
.PHONY: help
help:
	@echo "AudioBridge Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build the binary"
	@echo "  build-all     Build for all platforms"
	@echo "  lint          Run golangci-lint"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  clean         Remove build artifacts"
	@echo "  mod-tidy      Run go mod tidy"
	@echo "  deps          Download dependencies"
	@echo "  help          Show this help"
