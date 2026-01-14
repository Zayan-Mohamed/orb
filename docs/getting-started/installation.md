# Installation

This guide will help you install Orb on your system.

## Prerequisites

- Go 1.21 or higher (for building from source)
- Git (for cloning the repository)
- Internet connection (for downloading dependencies)

## Binary Releases

The easiest way to install Orb is to download a pre-built binary from the [releases page](https://github.com/Zayan-Mohamed/orb/releases).

### Linux

```bash
# Download the latest Linux binary
wget https://github.com/Zayan-Mohamed/orb/releases/latest/download/orb-linux-amd64

# Make it executable
chmod +x orb-linux-amd64

# Move to a directory in your PATH
sudo mv orb-linux-amd64 /usr/local/bin/orb

# Verify installation
orb --version
```

### macOS

```bash
# Download the latest macOS binary
curl -LO https://github.com/Zayan-Mohamed/orb/releases/latest/download/orb-darwin-amd64

# Make it executable
chmod +x orb-darwin-amd64

# Move to a directory in your PATH
sudo mv orb-darwin-amd64 /usr/local/bin/orb

# Verify installation
orb --version
```

For Apple Silicon (M1/M2):

```bash
curl -LO https://github.com/Zayan-Mohamed/orb/releases/latest/download/orb-darwin-arm64
chmod +x orb-darwin-arm64
sudo mv orb-darwin-arm64 /usr/local/bin/orb
```

### Windows

1. Download `orb-windows-amd64.exe` from the [releases page](https://github.com/Zayan-Mohamed/orb/releases)
2. Rename it to `orb.exe`
3. Add it to your PATH or place it in a directory that's already in your PATH
4. Open PowerShell or Command Prompt and run:
   ```
   orb --version
   ```

## Building from Source

If you prefer to build from source or want the latest development version:

### 1. Clone the Repository

```bash
git clone https://github.com/Zayan-Mohamed/orb.git
cd orb
```

### 2. Build with Make

```bash
# Build for your current platform
make build

# The binary will be created as ./orb
./orb --version
```

### 3. Build with Go

```bash
# Build for your current platform
go build -o orb

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o orb-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o orb-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o orb-windows-amd64.exe
```

### 4. Build with Build Script

```bash
# Build for all supported platforms
chmod +x build.sh
./build.sh

# Binaries will be created in ./build/ directory
```

### 5. Install Globally

```bash
# Install to $GOPATH/bin
go install

# Or manually copy to /usr/local/bin
sudo cp orb /usr/local/bin/
```

## Verifying Installation

After installation, verify that Orb is working correctly:

```bash
# Check version
orb --version

# Display help
orb --help

# List available commands
orb
```

## Dependencies

Orb requires the following Go modules (automatically downloaded during build):

- `github.com/spf13/cobra` - CLI framework
- `github.com/gorilla/websocket` - WebSocket implementation
- `golang.org/x/crypto` - Cryptography primitives
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - TUI components
- `github.com/charmbracelet/lipgloss` - TUI styling

## Next Steps

- Read the [First Steps](first-steps.md) guide to get started
- Check out [Usage Examples](examples.md) for common scenarios
- Learn about [Sharing Files](../user-guide/sharing.md) and [Connecting](../user-guide/connecting.md)

## Troubleshooting

### Permission Denied

If you get "permission denied" errors on Linux/macOS:

```bash
chmod +x orb
```

### Command Not Found

If `orb` command is not found:

- Make sure the directory is in your PATH
- Try running with full path: `/usr/local/bin/orb`
- Verify the binary is in the expected location

### Build Errors

If building from source fails:

- Ensure Go 1.21+ is installed: `go version`
- Clear the module cache: `go clean -modcache`
- Update dependencies: `go mod tidy`
