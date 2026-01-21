.PHONY: help build test clean example deps lint format cli cli-release

VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X github.com/techulus/capture-go/internal/cli.Version=$(VERSION) -X github.com/techulus/capture-go/internal/cli.Commit=$(COMMIT) -X github.com/techulus/capture-go/internal/cli.BuildDate=$(BUILD_DATE)"

help:
	@echo "Available targets:"
	@echo "  build       - Build the package"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  example     - Build and run example"
	@echo "  deps        - Download dependencies"
	@echo "  lint        - Run linter"
	@echo "  format      - Format code"
	@echo "  cli         - Build the CLI"
	@echo "  cli-release - Build CLI for all platforms"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build the package
build: deps
	go build -o bin/capture-go .

# Run tests
test: deps
	go test -v ./...

# Run tests with coverage
test-coverage: deps
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Build and run example
example: deps
	@echo "Building example..."
	go build -o bin/example ./example
	@echo "Example built. Run with: ./bin/example"

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Format code
format:
	go fmt ./...
	gofmt -s -w .

# Install dependencies for development
install-dev:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Generate documentation
docs:
	godoc -http=:6060

# Run benchmarks
bench: deps
	go test -bench=. ./...

# Check for security vulnerabilities
security:
	gosec ./...

# Create release build
release: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/capture-go-linux-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o bin/capture-go-darwin-amd64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -o bin/capture-go-windows-amd64.exe .

cli: deps
	go build $(LDFLAGS) -o bin/capture ./cmd/capture

cli-release: clean deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -a -installsuffix cgo -o bin/capture-linux-amd64 ./cmd/capture
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -a -installsuffix cgo -o bin/capture-linux-arm64 ./cmd/capture
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -a -installsuffix cgo -o bin/capture-darwin-amd64 ./cmd/capture
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -a -installsuffix cgo -o bin/capture-darwin-arm64 ./cmd/capture
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -a -installsuffix cgo -o bin/capture-windows-amd64.exe ./cmd/capture
