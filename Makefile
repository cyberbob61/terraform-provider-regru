# Provider's version
VERSION := 0.3.0

# Defining the path to install the plugin
OS_ARCH := $(shell go env GOOS)_$(shell go env GOARCH)
PLUGIN_DIR := /app/out
PLUGIN_NAME := terraform-provider-regru
BUILD_OUTPUT := $(PLUGIN_DIR)/$(PLUGIN_NAME)

# Compiling and installing the provider
.PHONY: build
build:
	@echo "Building Terraform provider..."
	mkdir -p $(PLUGIN_DIR)
	go build -o $(BUILD_OUTPUT)
	@echo "Build completed and installed at $(BUILD_OUTPUT)"

# Getting version Go из go.mod
.PHONY: go-version
go-version:
	@grep ^go go.mod | awk '{ print $$2 }'

# Running tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./... -v
	@echo "Tests completed"

# Formatting code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# Linting code
.PHONY: lint
lint:
	@echo "Linting code..."
	go vet ./...
	@echo "Linting completed"

# Cleaning up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BUILD_OUTPUT)
	@echo "Cleanup completed"

# Installing all dependencies
.PHONY: install-deps
install-deps:
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Dependencies installed"

# Help
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build          Build and install the Terraform provider"
	@echo "  make go-version Get the Go version from go.mod"
	@echo "  make test           Run tests"
	@echo "  make fmt            Format the code"
	@echo "  make lint           Lint the code"
	@echo "  make clean          Clean up build artifacts"
	@echo "  make install-deps   Install all dependencies"
	@echo "  make help           Display this help message"