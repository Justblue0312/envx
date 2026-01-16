.PHONY: help test test-race test-coverage benchmark lint clean build install release

# Default target
help:
	@echo "Available targets:"
	@echo "  test          - Run tests"
	@echo "  test-race     - Run tests with race detector"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  benchmark     - Run benchmarks"
	@echo "  lint          - Run linter"
	@echo "  clean         - Clean build artifacts"
	@echo "  build         - Build the package"
	@echo "  install       - Install the package"
	@echo "  release       - Run goreleaser (for releases only)"

# Testing
test:
	go test -v ./...

test-race:
	go test -race -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

benchmark:
	go test -bench=. -benchmem ./...

# Code quality
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Security
security:
	@which govulncheck > /dev/null || (echo "Installing govulncheck..." && go install golang.org/x/vuln/cmd/govulncheck@latest)
	govulncheck ./...

# Building
clean:
	rm -f coverage.out coverage.html
	rm -rf bin/
	rm -rf dist/

build:
	@mkdir -p bin
	go build -o bin/envx ./cmd/envx

build-lib:
	go build ./...

install:
	go install ./cmd/envx
	go install .

# Release (requires goreleaser)
release:
	@which goreleaser > /dev/null || (echo "Installing goreleaser..." && go install github.com/goreleaser/goreleaser@latest)
	goreleaser release --snapshot --clean

# Development setup
setup:
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/goreleaser/goreleaser@latest

# CI commands
ci: test-race lint security