# Building from Source

Instructions for building Orb from source code.

## Prerequisites

- Go 1.21 or higher
- Git
- Make (optional)

## Clone Repository

```bash
git clone https://github.com/Zayan-Mohamed/orb.git
cd orb
```

## Build Methods

### Method 1: Using Make

```bash
# Build for current platform
make build

# Run tests
make test

# Clean build artifacts
make clean
```

### Method 2: Using Go

```bash
# Build
go build -o orb

# Build with optimizations
go build -ldflags="-s -w" -o orb

# Install to $GOPATH/bin
go install
```

### Method 3: Using Build Script

```bash
# Build for all platforms
chmod +x build.sh
./build.sh

# Output in build/ directory
ls build/
```

## Cross-Compilation

### Linux

```bash
GOOS=linux GOARCH=amd64 go build -o orb-linux-amd64
GOOS=linux GOARCH=arm64 go build -o orb-linux-arm64
```

### macOS

```bash
GOOS=darwin GOARCH=amd64 go build -o orb-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o orb-darwin-arm64
```

### Windows

```bash
GOOS=windows GOARCH=amd64 go build -o orb-windows-amd64.exe
```

## Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...

# Specific package
go test ./internal/crypto

# Run benchmarks
go test -bench=. ./...
```

## Development Workflow

```bash
# Install dependencies
go mod download

# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linters
golangci-lint run

# Build and run
go run main.go relay
```

## Build Configuration

### Build Tags

```bash
# Build with debug symbols
go build -tags debug

# Build without CGO
CGO_ENABLED=0 go build
```

### Optimization

```bash
# Strip symbols
go build -ldflags="-s -w"

# Set version
go build -ldflags="-X main.version=1.0.0"
```

## Troubleshooting

### Module Issues

```bash
go mod tidy
go clean -modcache
```

### Build Errors

```bash
# Verbose build
go build -x -v

# Check Go version
go version
```

## Next Steps

- [Contributing Guidelines](contributing.md)
- [API Reference](api.md)
