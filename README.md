# Orb — Zero-Trust Folder Tunneling Tool

[![Build Status](https://github.com/Zayan-Mohamed/orb/workflows/Build/badge.svg)](https://github.com/Zayan-Mohamed/orb/actions/workflows/build.yml)
[![Release](https://github.com/Zayan-Mohamed/orb/workflows/Release/badge.svg)](https://github.com/Zayan-Mohamed/orb/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Zayan-Mohamed/orb)](https://goreportcard.com/report/github.com/Zayan-Mohamed/orb)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Zayan-Mohamed/orb)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-mkdocs-blue.svg)](https://zayan-mohamed.github.io/orb/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](docs/development/contributing.md)

Orb is a secure, terminal-first utility that allows you to share a local folder across the internet using end-to-end encryption. No accounts, no cloud storage, no port forwarding required.

## Demo

### Sharing a Folder

<!-- TODO: Add GIF recording of 'orb share' command -->

![Orb Share Demo](docs/assets/images/orb-share-demo.gif)
_Share a folder with a single command - encrypted end-to-end_

### Browsing Shared Files

<!-- TODO: Add GIF recording of TUI file browser -->

![Orb Connect Demo](docs/assets/images/orb-connect-demo.gif)
_Interactive TUI browser for secure file access_

<!-- ### Quick Transfer -->

<!-- TODO: Add GIF showing complete workflow from share to download -->

<!-- ![Orb Complete Demo](docs/assets/images/orb-complete-demo.gif)
_Complete workflow: Share → Connect → Browse → Download_ -->

## Features

- **Zero-Trust Architecture**: The relay server never sees plaintext data
- **Strong Cryptography**: Argon2id for key derivation, Noise Protocol for handshake, ChaCha20-Poly1305 for transport encryption
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **NAT-Safe**: All connections are outbound, works behind firewalls
- **TUI File Browser**: Interactive terminal interface for browsing and downloading files
- **No Long-Term Secrets**: Sessions expire automatically
- **Secure by Design**: Path sanitization, symlink protection, replay protection

## Quick Start

### Install

```bash
# Download the binary for your platform
curl -L https://github.com/Zayan-Mohamed/orb/releases/latest/download/orb-$(uname -s)-$(uname -m) -o orb
chmod +x orb
sudo mv orb /usr/local/bin/
```

Or [build from source](#building-from-source).

### Share a Folder

```bash
orb share ./myfolder
```

Output:

```
Session ID: abc123def456
Passcode: secure-random-passcode
Relay: ws://localhost:8080

Share these credentials securely with the recipient.
Waiting for connection...
```

### Connect to a Shared Folder

```bash
orb connect 7F9Q2A
```

Prompts for passcode, then opens an interactive file browser.

## Commands

### `orb share <path>`

Share a local directory over an encrypted tunnel.

Options:

- `--relay <url>`: Relay server URL (default: http://localhost:8080)
- `--readonly`: Share folder in read-only mode

Example:

```bash
orb share ./documents --readonly
```

### `orb connect <session-id>`

Connect to a shared session.

Options:

- `--relay <url>`: Relay server URL
- `--passcode <code>`: Session passcode (prompts if not provided)
- `--tui`: Use TUI file browser (default: true)
- `--mount <path>`: Mount point for FUSE (Linux/macOS only)

Example:

```bash
orb connect 7F9Q2A --passcode 493-771 --tui
```

### `orb relay`

Start a relay server.

Options:

- `--listen <addr>`: Listen address (default: :8080)

Example:

```bash
orb relay --listen 0.0.0.0:8080
```

## Security Features

### Cryptography

- **Key Derivation**: Argon2id with 64MB memory, 3 iterations
- **Handshake**: Noise Protocol Framework with X25519 key exchange
- **Transport**: ChaCha20-Poly1305 authenticated encryption
- **Random Generation**: Cryptographically secure random numbers

### Protection Against

- Passcode brute force (Argon2id memory-hard function)
- Replay attacks (unique nonces per packet)
- Man-in-the-middle (Noise Protocol mutual authentication)
- Path traversal (path sanitization and validation)
- Symlink attacks (symlink resolution and boundary checking)
- Session hijacking (rate limiting and session locking)
- Relay compromise (end-to-end encryption, relay is blind)

### Filesystem Security

- All paths are sanitized and validated
- Symlinks pointing outside the shared directory are blocked
- No execution of files remotely
- Configurable read-only mode
- Automatic session expiration

## Architecture

Orb consists of three components:

1. **CLI**: User interface for sharing and connecting
2. **Relay Server**: Blind byte pipe that forwards encrypted data
3. **Encrypted Tunnel**: End-to-end encrypted communication channel

```
┌─────────┐                  ┌───────────┐                  ┌─────────┐
│ Sharer  │◄────encrypted────►│   Relay   │◄────encrypted────►│Receiver │
│  (CLI)  │                  │  Server   │                  │  (CLI)  │
└─────────┘                  └───────────┘                  └─────────┘
     │                              │                              │
     │                              │                              │
     └──────────────────────────────┴──────────────────────────────┘
              All encryption happens at the edges
              Relay never sees plaintext data
```

## Building from Source

### Prerequisites

- Go 1.22 or later
- Make (optional, but recommended)

### Build

```bash
# Clone the repository
git clone https://github.com/Zayan-Mohamed/orb.git
cd orb

# Install dependencies
go mod download

# Build for current platform
make build-local

# Or build for all platforms
./build.sh
```

### Cross-Platform Builds

```bash
# Build for all platforms
./build.sh

# Binaries will be in build/
ls build/
# orb-linux-amd64
# orb-linux-arm64
# orb-darwin-amd64
# orb-darwin-arm64
# orb-windows-amd64.exe
```

## Development

### Run Tests

```bash
make test
```

### Code Coverage

```bash
go test -cover ./...
```

### Run Relay Server Locally

```bash
make relay
```

### Run Sharer

```bash
make share
```

### Run Receiver

```bash
SESSION=<session-id> make connect
```

## Configuration

Orb uses sensible defaults and requires no configuration files. All settings are passed via command-line flags.

## Documentation

For comprehensive documentation, visit the [Orb Documentation](docs/):

- **Getting Started**

  - [Installation Guide](docs/getting-started/installation.md)
  - [First Steps](docs/getting-started/first-steps.md)
  - [Usage Examples](docs/getting-started/examples.md)

- **User Guides**

  - [Sharing Files](docs/user-guide/sharing.md)
  - [Connecting](docs/user-guide/connecting.md)
  - [TUI Browser](docs/user-guide/tui.md)
  - [Troubleshooting](docs/user-guide/troubleshooting.md)

- **Security**

  - [Cryptography Details](docs/security/cryptography.md)
  - [Best Practices](docs/security/best-practices.md)
  - [Threat Model](docs/security/threat-model.md)

- **Deployment**

  - [Relay Server](docs/deployment/relay-server.md)
  - [Production](docs/deployment/production.md)
  - [Docker](docs/deployment/docker.md)

- **Development**
  - [Building from Source](docs/development/building.md)
  - [Contributing](docs/development/contributing.md)
  - [API Reference](docs/development/api.md)

### Building Documentation Site

```bash
# Install MkDocs
pip install -r requirements.txt

# Serve locally
mkdocs serve

# Build static site
mkdocs build
```

## Platform-Specific Notes

### Linux

- TUI mode works out of the box
- Install: Use provided binaries or build from source

### macOS

- TUI mode works out of the box
- Install: Use provided binaries or build from source

### Windows

- TUI mode works out of the box
- Install: Use provided binaries or build from source

## Limitations

- Maximum file size: Limited by available memory (10MB per read operation)
- Sessions expire after 24 hours
- Maximum 5 failed passcode attempts before session lock
- Read-only access (no file upload/modification)

## Security Disclosure

If you discover a security vulnerability, please email security@orb.example.com. Do not create public issues for security vulnerabilities.

## License

MIT License - see [LICENSE](docs/about/license.md) file for details

## Contributing

Contributions are welcome! Please read [Contributing Guide](docs/development/contributing.md) for guidelines.

## Acknowledgments

- [Noise Protocol Framework](https://noiseprotocol.org/)
- [ChaCha20-Poly1305 (RFC 8439)](https://www.rfc-editor.org/rfc/rfc8439)
- [Argon2id (RFC 9106)](https://www.rfc-editor.org/rfc/rfc9106)
- [Bubble Tea TUI framework](https://github.com/charmbracelet/bubbletea)
- [Go crypto libraries](https://golang.org/x/crypto)

---

<div align="center">

**Remember**: Orb is designed for security. If a feature weakens encryption, privacy, or isolation, it will not be implemented.

**Stay safe. Stay secure.**

Made with care by developers who prioritize security and privacy.

[Star us on GitHub](https://github.com/Zayan-Mohamed/orb) | [Read the Docs](https://zayan-mohamed.github.io/orb/) | [Report Bug](https://github.com/Zayan-Mohamed/orb/issues) | [Request Feature](https://github.com/Zayan-Mohamed/orb/issues)

</div>
