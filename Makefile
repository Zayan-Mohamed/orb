.PHONY: build clean test relay share connect deps

# Build all binaries
build:
	@./build.sh

# Build single binary for current platform
build-local:
	go build -o orb .

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run relay server
relay:
	go run . relay --listen :8080

# Share a folder (example)
share:
	go run . share ./test --relay http://localhost:8080

# Connect to a session (example)
connect:
	go run . connect $(SESSION) --relay http://localhost:8080

# Run tests
test:
	go test -v ./...

# Run security tests
test-security:
	go test -v -run Security ./...

# Clean build artifacts
clean:
	rm -rf build/
	rm -f orb

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install locally
install:
	go install .

# Cross-compile for all platforms
release:
	VERSION=$(shell git describe --tags --always --dirty) ./build.sh
