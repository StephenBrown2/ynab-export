# Justfile for ynab-export
# Requires: https://github.com/casey/just
# Export GOEXPERIMENT to enable json/v2 for all recipes

export GOEXPERIMENT := "jsonv2"

# Use Git Bash on Windows for better compatibility

set windows-shell := ["bash.exe", "-c"]

# Binary name based on OS

bin_file := if os_family() == "windows" { "ynab-export.exe" } else { "ynab-export" }

# Default recipe to display help information
default:
    @just --list

# Build the application for the current platform
[group('build')]
build:
    go build -o {{ bin_file }} .

[group('build')]
install:
    go install .

# Run the application
[group('dev')]
run:
    go run .

# Run tests
[group('test')]
test:
    go test -v ./...

# Clean build artifacts
[group('build')]
clean:
    rm -f {{ bin_file }}
    rm -rf dist/
    go clean

# Install dependencies and golangci-lint
[group('setup')]
deps:
    go mod download
    go mod tidy
    @echo "Installing golangci-lint..."
    @curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Cross-compile for all platforms (both amd64 and arm64)
[group('build')]
build-all: build-linux build-linux-arm build-darwin build-darwin-arm build-windows build-windows-arm

# Build for Linux (amd64)
[group('build')]
build-linux:
    GOOS=linux GOARCH=amd64 go build -o dist/ynab-export-linux-amd64 .

# Build for Linux (arm64)
[group('build')]
build-linux-arm:
    GOOS=linux GOARCH=arm64 go build -o dist/ynab-export-linux-arm64 .

# Build for macOS (amd64 - Intel)
[group('build')]
build-darwin:
    GOOS=darwin GOARCH=amd64 go build -o dist/ynab-export-darwin-amd64 .

# Build for macOS (arm64 - Apple Silicon)
[group('build')]
build-darwin-arm:
    GOOS=darwin GOARCH=arm64 go build -o dist/ynab-export-darwin-arm64 .

# Build for Windows (amd64)
[group('build')]
build-windows:
    GOOS=windows GOARCH=amd64 go build -o dist/ynab-export-windows-amd64.exe .

# Build for Windows (arm64)
[group('build')]
build-windows-arm:
    GOOS=windows GOARCH=arm64 go build -o dist/ynab-export-windows-arm64.exe .

# Format code using golangci-lint formatters
[group('lint')]
fmt:
    just --fmt --unstable
    @echo "Formatting code..."
    @if command -v golangci-lint >/dev/null 2>&1; then \
        golangci-lint fmt ./...; \
    else \
        echo "ERROR: golangci-lint not found. Run 'just deps' to install it."; \
        exit 1; \
    fi

# Check formatting without making changes
[group('lint')]
fmt-check:
    #!/usr/bin/env bash
    set -euo pipefail
    just --fmt --check --unstable
    echo "Checking code formatting..."
    if command -v golangci-lint >/dev/null 2>&1; then
        output=$(golangci-lint fmt --diff ./... 2>&1)
        if [ -n "$output" ]; then
            echo "$output"
            echo "Files need formatting. Run 'just fmt' to fix."
            exit 1
        else
            echo "All files are formatted correctly"
        fi
    else
        echo "ERROR: golangci-lint not found. Run 'just deps' to install it."
        exit 1
    fi

# Run linter with golangci-lint
[group('lint')]
lint:
    @echo "Running golangci-lint..."
    @if command -v golangci-lint >/dev/null 2>&1; then \
        golangci-lint run ./... || { \
            echo ""; \
            echo "Linting failed! Try running 'just lint-fix' first to auto-fix issues."; \
            echo "After that, manually address any remaining issues."; \
            exit 1; \
        }; \
    else \
        echo "ERROR: golangci-lint not found. Run 'just deps' to install it."; \
        exit 1; \
    fi

# Run linter and automatically fix issues where possible
[group('lint')]
lint-fix: fmt
    @echo "Running golangci-lint with auto-fix..."
    @golangci-lint run --fix ./...

# Run all checks (format, lint, test)
[group('test')]
check: fmt-check lint test
    @echo "All checks passed!"

# Show Go environment information
[group('setup')]
info:
    @echo "Go version:"
    @go version
    @echo ""
    @echo "Go environment:"
    @go env GOOS GOARCH
    @echo ""
    @echo "Module info:"
    @go list -m

# Create a new release (tags and pushes)
[group('release')]
release version:
    @echo "Creating release {{ version }}"
    @git push origin master
    git tag -a {{ version }} -m "Release {{ version }}"
    git push origin {{ version }}
    @echo
    @echo "Monitor the release on GitHub:"
    @echo "https://github.com/StephenBrown2/ynab-export/actions"
    @echo "https://github.com/StephenBrown2/ynab-export/releases/tag/{{ version }}"

# Delete a release tag and recreate it (useful for fixing failed releases)
[group('release')]
redo-release version:
    @echo "Deleting tag {{ version }} locally and remotely..."
    -git tag -d {{ version }}
    -git push origin :refs/tags/{{ version }}
    @echo "Recreating release {{ version }}"
    @just release {{ version }}
