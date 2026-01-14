# API Reference

Internal API documentation for Orb developers.

## Package Structure

```
orb/
├── cmd/                # CLI commands
├── internal/           # Internal packages
│   ├── crypto/        # Cryptography
│   ├── filesystem/    # File operations
│   ├── relay/         # Relay server
│   ├── session/       # Session management
│   ├── tui/           # Terminal UI
│   └── tunnel/        # Encrypted tunnel
├── pkg/               # Public packages
│   └── protocol/      # Wire protocol
└── main.go            # Entry point
```

## Key Packages

### crypto

**Purpose:** Cryptographic primitives

**Types:**

```go
type AEAD struct {
    cipher    cipher.AEAD
    sendNonce uint64
    recvNonce uint64
}

func NewAEAD(key []byte) (*AEAD, error)
func (a *AEAD) Encrypt(plaintext []byte) ([]byte, error)
func (a *AEAD) Decrypt(ciphertext []byte) ([]byte, error)
```

**Functions:**

```go
func DeriveKey(passcode, sessionID []byte) []byte
func GenerateKeyPair() (private, public []byte, err error)
```

### filesystem

**Purpose:** Secure file operations

**Functions:**

```go
func SanitizePath(base, requested string) (string, error)
func ReadFile(base, path string) ([]byte, error)
func ListDirectory(base, path string) ([]FileInfo, error)
```

### tunnel

**Purpose:** Encrypted communication channel

**Types:**

```go
type Tunnel struct {
    conn   *websocket.Conn
    aead   *crypto.AEAD
}

func NewTunnel(conn *websocket.Conn, key []byte) *Tunnel
func (t *Tunnel) Send(frame *protocol.Frame) error
func (t *Tunnel) Receive() (*protocol.Frame, error)
```

### protocol

**Purpose:** Wire protocol definitions

**Types:**

```go
type FrameType uint8

const (
    FrameTypeRequest FrameType = iota
    FrameTypeResponse
    FrameTypeError
)

type Frame struct {
    Type      FrameType
    Operation string
    Path      string
    Data      []byte
}
```

## CLI Commands

### share

```go
func shareCmd() *cobra.Command {
    // Implementation in cmd/share.go
}
```

### connect

```go
func connectCmd() *cobra.Command {
    // Implementation in cmd/connect.go
}
```

### relay

```go
func relayCmd() *cobra.Command {
    // Implementation in cmd/relay.go
}
```

## Configuration

Currently no configuration files. All options via flags.

## Error Handling

Use standard Go error handling:

```go
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

## Testing

```bash
go test ./...
```

## Next Steps

- [Building from Source](building.md)
- [Contributing](contributing.md)
