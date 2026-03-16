MOCKGEN := mockgen
REPO_INTERFACES_DIR := internal/repository/interfaces
REPO_MOCK_DIR := internal/repository/mock

# Detect OS and Architecture
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
BINARY_NAME := server
BINARY_PATH := ./bin/$(BINARY_NAME)-$(GOOS)-$(GOARCH)
ifeq ($(GOOS),windows)
	BINARY_PATH := ./bin/$(BINARY_NAME)-$(GOOS)-$(GOARCH).exe
endif

.PHONY: help dev server build test install-deps clean repository-mocks

help:
	@echo "Available commands:"
	@echo "  make install-deps  - Install Go dependencies"
	@echo "  make dev           - Start server with hot-reload"
	@echo "  make build         - Build production binary"
	@echo "  make test          - Run all tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Clean build artifacts"

install-deps:
	go mod tidy
	go mod download

dev:
	DEV_MODE=true air

server:
	DEV_MODE=true go run ./cmd/server

build:
	@echo "Building server binary for $(GOOS)/$(GOARCH)..."
	@mkdir -p ./bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_PATH) ./cmd/server/main.go
	@echo "Binary created at: $(BINARY_PATH)"

test:
	go test -v -cover -race ./...

test-verbose:
	go test -v -cover -race -failfast ./...

test-coverage:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf ./bin
	rm -f coverage.out coverage.html
	go clean -testcache

repository-mocks:
	@rm -rf $(REPO_MOCK_DIR)
	@mkdir -p $(REPO_MOCK_DIR)
	@for f in $(REPO_INTERFACES_DIR)/*.go; do \
		base=$$(basename $$f .go); \
		name=$${base%_interface}; \
		echo "Generating mock for interface: $$name"; \
		$(MOCKGEN) \
			-source=$$f \
			-destination=$(REPO_MOCK_DIR)/$${name}_mock.go \
			-package=mock; \
	done
