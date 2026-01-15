# Orb Architecture

## Overview

Orb is a zero-trust file sharing system that enables secure, peer-to-peer folder access over untrusted networks. The architecture is designed around three core principles:

1. **Zero-Trust Model**: The relay server never has access to plaintext data or decryption keys
2. **End-to-End Encryption**: All cryptographic operations occur at the endpoints
3. **Minimal Attack Surface**: Simple, auditable components with clear security boundaries

This document describes the system architecture, component interactions, and key design decisions.

## System Architecture

### High-Level Components

Orb consists of three primary components that work together to enable secure file sharing:

```
┌─────────────────────────────────────────────────────────────────┐
│                         Orb System                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│  │   Sharer     │         │    Relay     │         │  Receiver    │
│  │   (Client)   │◄───────►│   Server     │◄───────►│   (Client)   │
│  └──────────────┘         └──────────────┘         └──────────────┘
│         │                        │                         │
│         │                        │                         │
│    ┌────▼─────┐            ┌────▼─────┐            ┌────▼─────┐
│    │  Crypto  │            │  WebSocket│            │  Crypto  │
│    │  Layer   │            │   Router  │            │  Layer   │
│    └────┬─────┘            └──────────┘            └────┬─────┘
│         │                                                │
│    ┌────▼─────┐                                    ┌────▼─────┐
│    │   File   │                                    │   TUI    │
│    │  System  │                                    │  Browser │
│    └──────────┘                                    └──────────┘
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### 1. Sharer (Server-side Client)

The Sharer is responsible for:

- Exposing a local directory for secure access
- Generating session credentials (ID and passcode)
- Performing filesystem operations in a sandboxed environment
- Encrypting all outbound data

#### 2. Relay Server

The Relay Server acts as a blind intermediary that:

- Routes encrypted traffic between peers
- Manages WebSocket connections
- Pairs clients based on session IDs
- Has zero knowledge of plaintext data or keys

**Critical Security Property**: The relay server cannot decrypt any traffic. It only sees encrypted bytes.

#### 3. Receiver (Client-side)

The Receiver enables users to:

- Connect to a shared session using credentials
- Browse the remote filesystem through a TUI
- Download files over the encrypted tunnel
- Verify the authenticity of the connection

## Connection Flow

### Session Establishment

```
┌─────────┐                    ┌───────────┐                    ┌─────────┐
│ Sharer  │                    │   Relay   │                    │Receiver │
└────┬────┘                    └─────┬─────┘                    └────┬────┘
     │                               │                               │
     │ 1. Create Session             │                               │
     │──────────────────────────────►│                               │
     │ Session ID + Passcode         │                               │
     │◄──────────────────────────────│                               │
     │                               │                               │
     │ 2. Connect to Relay           │                               │
     │──────────────────────────────►│                               │
     │ WebSocket established         │                               │
     │                               │                               │
     │                               │ 3. Connect with Session ID    │
     │                               │◄──────────────────────────────│
     │                               │ WebSocket established         │
     │                               │                               │
     │                               │ 4. Pair connections           │
     │                               │───────────┐                   │
     │                               │           │                   │
     │                               │◄──────────┘                   │
     │                               │                               │
     │ 5. Noise Handshake (encrypted through relay)                 │
     │◄──────────────────────────────┼──────────────────────────────►│
     │    - Exchange ephemeral keys  │                               │
     │    - Derive shared secret     │                               │
     │    - Mutual authentication    │                               │
     │                               │                               │
     │ 6. Encrypted Tunnel Established                               │
     │◄══════════════════════════════╪═══════════════════════════════►│
     │    All data encrypted E2E     │   Relay blind forwards        │
     │                               │                               │
```

### Phase-by-Phase Breakdown

#### Phase 1: Session Creation

1. Sharer requests session from relay server
2. Relay generates unique session ID
3. Sharer derives encryption keys from user-provided passcode using Argon2id
4. Session credentials (ID + passcode) shared out-of-band with receiver

#### Phase 2: Connection Pairing

1. Both sharer and receiver connect to relay via WebSocket
2. Relay pairs connections with matching session ID
3. Relay begins forwarding encrypted frames between peers

#### Phase 3: Cryptographic Handshake

1. Peers perform Noise Protocol handshake through relay
2. Ephemeral X25519 keypairs generated on each side
3. Mutual authentication using passcode-derived keys
4. Shared transport keys established with perfect forward secrecy

#### Phase 4: Encrypted Communication

1. All subsequent messages encrypted with ChaCha20-Poly1305
2. Each frame has unique nonce for replay protection
3. Filesystem operations (LIST, READ, etc.) transmitted as encrypted protocol frames

## Core Subsystems

### 1. Cryptographic Layer (`internal/crypto/`)

The cryptographic layer provides all security primitives and is built on industry-standard algorithms.

#### Key Derivation

```go
// Passcode → Encryption Keys
func DeriveKey(passcode, sessionID []byte) []byte {
    return argon2.IDKey(
        passcode,           // User input
        sessionID,          // Salt (unique per session)
        3,                  // Iterations
        64*1024,           // Memory (64 MB)
        4,                  // Parallelism
        32,                 // Output length
    )
}
```

**Design Rationale**: Argon2id is memory-hard, making brute-force attacks expensive even with specialized hardware.

#### Noise Protocol Implementation

Orb uses a simplified Noise XX pattern for the handshake:

```
Initiator                   Responder
──────────────────────────────────────────
Generate ephemeral key
  e (X25519)
                            Generate ephemeral key
                              e' (X25519)

Send: e, encrypted auth
─────────────────────────►
                            Verify auth
                            Send: e', encrypted auth
                            ◄─────────────────────────
Verify auth

Derive transport keys:      Derive transport keys:
  tk_send = KDF(e, e')        tk_recv = KDF(e, e')
  tk_recv = KDF(e, e')        tk_send = KDF(e, e')
```

**Security Properties**:

- **Perfect Forward Secrecy**: Session keys destroyed after handshake
- **Mutual Authentication**: Both parties prove knowledge of passcode
- **Identity Hiding**: No static identities transmitted

#### Transport Encryption

After handshake, all messages encrypted with ChaCha20-Poly1305 AEAD:

```
Plaintext Frame → [Nonce (24B)][Ciphertext][Auth Tag (16B)]
```

**Nonce Management**: Counter-based nonces prevent reuse and provide replay protection.

### 2. Protocol Layer (`pkg/protocol/`)

The wire protocol defines binary frames for all filesystem operations.

#### Frame Structure

```
┌────────────────────────────────────────────────────────────┐
│                     Protocol Frame                         │
├─────────────┬──────────────┬──────────────┬────────────────┤
│  Type (1B)  │  Length (4B) │  Payload (N) │  Reserved (3B) │
└─────────────┴──────────────┴──────────────┴────────────────┘
```

#### Operation Types

| Type | Operation | Direction | Description             |
| ---- | --------- | --------- | ----------------------- |
| 0x01 | LIST      | R→S       | List directory contents |
| 0x02 | STAT      | R→S       | Get file metadata       |
| 0x03 | READ      | R→S       | Read file contents      |
| 0x04 | WRITE     | R→S       | Write file contents     |
| 0x05 | DELETE    | R→S       | Delete file             |
| 0x06 | MKDIR     | R→S       | Create directory        |
| 0x07 | RENAME    | R→S       | Rename file/directory   |
| 0x10 | RESPONSE  | S→R       | Operation response      |
| 0xFF | ERROR     | S→R       | Error response          |

**R**: Receiver (client), **S**: Sharer (server)

#### Example: LIST Operation

```
Request (Receiver → Sharer):
┌──────┬──────────┬──────────────────┐
│ 0x01 │ 0x000005 │ "/docs"          │
└──────┴──────────┴──────────────────┘

Response (Sharer → Receiver):
┌──────┬──────────┬────────────────────────────────────────┐
│ 0x10 │ 0x000042 │ ["file1.txt", "file2.pdf", "subdir/"] │
└──────┴──────────┴────────────────────────────────────────┘
```

### 3. Secure Filesystem (`internal/filesystem/`)

The filesystem layer enforces security boundaries and prevents unauthorized access.

#### Path Sanitization Pipeline

```
User Input → Clean → Resolve → Validate → Safe Path
    ↓          ↓        ↓         ↓          ↓
"/../../etc"  "/etc"  "/etc"   [REJECT]    N/A
"docs/../"    "."     "/share" [ACCEPT]    "/share"
"file.txt"    "file"  "/share/file" [ACCEPT] "/share/file.txt"
```

**Validation Rules**:

1. All paths must resolve within the shared directory root
2. Symlinks pointing outside root are rejected
3. Attempts to escape via `..` are blocked
4. No absolute paths allowed in user input

#### Sandboxing Architecture

```
┌────────────────────────────────────────┐
│         Application Layer              │
└────────────────┬───────────────────────┘
                 │
                 ▼
┌────────────────────────────────────────┐
│      Secure Filesystem Wrapper         │
│  ┌──────────────────────────────────┐  │
│  │  Path Sanitization               │  │
│  └──────────────────────────────────┘  │
│  ┌──────────────────────────────────┐  │
│  │  Symlink Resolution & Validation │  │
│  └──────────────────────────────────┘  │
│  ┌──────────────────────────────────┐  │
│  │  Boundary Checking               │  │
│  └──────────────────────────────────┘  │
└────────────────┬───────────────────────┘
                 │
                 ▼
┌────────────────────────────────────────┐
│         Operating System               │
│         Real Filesystem                │
└────────────────────────────────────────┘
```

### 4. Session Management (`internal/session/`)

Sessions provide temporary, time-bound access with built-in security controls.

#### Session Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│                    Session Lifecycle                        │
└─────────────────────────────────────────────────────────────┘

Create → Active → [ Success / Locked / Expired ] → Destroyed
  │       │              │         │        │           │
  │       │              │         │        │           │
  │       └─ Processing  │         │        │           │
  │          requests    │         │        │           │
  │                      │         │        │           │
  │              After 5 failed    │        │           │
  │              auth attempts     │        │           │
  │                      ▼         │        │           │
  │                   LOCKED ──────┘        │           │
  │                                         │           │
  │                          After 24 hours │           │
  │                                         ▼           │
  └─────────────────────────────── EXPIRED ────────────┘
```

**Session Properties**:

- **Unique ID**: 12-character alphanumeric identifier
- **Ephemeral**: Automatically expire after 24 hours
- **Rate Limited**: Lock after 5 failed authentication attempts
- **Single-Use Pairing**: One receiver per session (current implementation)

### 5. Relay Server (`internal/relay/`)

The relay server is intentionally minimal to reduce attack surface and maintain zero-knowledge properties.

#### Relay Responsibilities

**What the Relay Does**:

- Accept WebSocket connections from clients
- Pair connections with matching session IDs
- Forward encrypted frames bidirectionally
- Manage connection lifecycle and timeouts

**What the Relay Cannot Do**:

- Decrypt any traffic (no keys)
- Authenticate users (no credentials)
- Inspect file operations (all encrypted)
- Log sensitive data (blind forwarding)

#### Connection Routing

```go
// Simplified relay logic
type Relay struct {
    sessions map[string]*SessionPair
}

func (r *Relay) HandleConnection(conn *websocket.Conn, sessionID string) {
    pair := r.sessions[sessionID]

    if pair.IsComplete() {
        // Route encrypted frames between paired connections
        go forwardFrames(pair.ClientA, pair.ClientB)
        go forwardFrames(pair.ClientB, pair.ClientA)
    } else {
        // Wait for second peer
        pair.AddPeer(conn)
    }
}

func forwardFrames(src, dst *websocket.Conn) {
    for {
        // Read encrypted frame from source
        msgType, data, err := src.ReadMessage()

        // Forward to destination (blind)
        dst.WriteMessage(msgType, data)
    }
}
```

### 6. Terminal UI (`internal/tui/`)

The TUI provides an interactive interface for browsing and downloading files.

#### Component Architecture

Built using the Bubble Tea framework with Model-View-Update pattern:

```
┌─────────────────────────────────────────┐
│              TUI Browser                │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────────────────────────────┐   │
│  │      View (Rendering)           │   │
│  │  - File list                    │   │
│  │  - Breadcrumbs                  │   │
│  │  - Status bar                   │   │
│  └─────────────────────────────────┘   │
│              │                          │
│              ▼                          │
│  ┌─────────────────────────────────┐   │
│  │      Model (State)              │   │
│  │  - Current directory            │   │
│  │  - File list                    │   │
│  │  - Selected item                │   │
│  └─────────────────────────────────┘   │
│              │                          │
│              ▼                          │
│  ┌─────────────────────────────────┐   │
│  │      Update (Logic)             │   │
│  │  - Handle keypresses            │   │
│  │  - Navigation                   │   │
│  │  - File operations              │   │
│  └─────────────────────────────────┘   │
│              │                          │
│              ▼                          │
│  ┌─────────────────────────────────┐   │
│  │   Commands (Side Effects)       │   │
│  │  - Fetch directory              │   │
│  │  - Download file                │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

## Code Organization

### Project Structure

```
orb/
├── cmd/                      # CLI entry points
│   ├── root.go              # Root command and global flags
│   ├── share.go             # Share command implementation
│   ├── connect.go           # Connect command implementation
│   ├── relay.go             # Relay server command
│   └── utils.go             # Shared CLI utilities
│
├── internal/                # Private application code
│   ├── crypto/
│   │   ├── crypto.go        # Key derivation, encryption primitives
│   │   └── noise.go         # Noise Protocol handshake
│   │
│   ├── filesystem/
│   │   └── secure_fs.go     # Sandboxed filesystem operations
│   │
│   ├── relay/
│   │   └── server.go        # WebSocket relay implementation
│   │
│   ├── session/
│   │   └── session.go       # Session lifecycle management
│   │
│   ├── tui/
│   │   └── browser.go       # Terminal user interface
│   │
│   └── tunnel/
│       └── tunnel.go        # Encrypted tunnel layer
│
├── pkg/                     # Public, reusable packages
│   └── protocol/
│       └── protocol.go      # Wire protocol definitions
│
├── main.go                  # Application entry point
└── go.mod                   # Go module dependencies
```

### Module Boundaries

```
┌─────────────────────────────────────────────────────────┐
│                     Public API (pkg/)                   │
│  - Protocol definitions                                 │
│  - Reusable across implementations                      │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                Internal Packages (internal/)            │
│  - Application-specific logic                           │
│  - Not importable by external code                      │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                    Commands (cmd/)                      │
│  - User-facing CLI interface                            │
│  - Orchestrates internal packages                       │
└─────────────────────────────────────────────────────────┘
```

## Security Architecture

### Defense in Depth

Orb implements multiple layers of security controls:

```
┌─────────────────────────────────────────────────────────┐
│ Layer 1: Network Security (Transport)                  │
│ - TLS for relay connections (wss://)                   │
│ - NAT traversal (no port forwarding needed)            │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│ Layer 2: Cryptographic Security                        │
│ - ChaCha20-Poly1305 AEAD encryption                    │
│ - Unique nonces (replay protection)                    │
│ - Perfect forward secrecy (ephemeral keys)             │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│ Layer 3: Authentication                                 │
│ - Argon2id key derivation (brute-force resistant)      │
│ - Mutual authentication (Noise Protocol)               │
│ - Rate limiting (5 attempts → lockout)                 │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│ Layer 4: Application Security                          │
│ - Path traversal prevention                            │
│ - Symlink escape protection                            │
│ - Resource limits (file size, frame size)              │
└─────────────────────────────────────────────────────────┘
```

### Threat Model & Mitigations

| Threat                | Attack Vector        | Mitigation                                          |
| --------------------- | -------------------- | --------------------------------------------------- |
| **Eavesdropping**     | Network monitoring   | End-to-end encryption (ChaCha20-Poly1305)           |
| **MITM**              | Network interception | Mutual authentication via Noise Protocol            |
| **Brute Force**       | Passcode guessing    | Argon2id (64MB, ~100ms per attempt) + rate limiting |
| **Replay**            | Packet resubmission  | Unique nonces, counter-based validation             |
| **Path Traversal**    | `../../etc/passwd`   | Path sanitization, root boundary enforcement        |
| **Symlink Escape**    | Symlink → `/etc`     | Symlink resolution with boundary checking           |
| **Session Hijacking** | Credential theft     | Rate limiting, session expiration (24h)             |
| **Relay Compromise**  | Server takeover      | Zero-knowledge design (no plaintext access)         |
| **DoS**               | Resource exhaustion  | Frame size limits, connection timeouts              |

### Cryptographic Details

#### Algorithms & Parameters

| Component          | Algorithm     | Parameters                            |
| ------------------ | ------------- | ------------------------------------- |
| **Key Derivation** | Argon2id      | 64 MB memory, 3 iterations, 4 threads |
| **Key Exchange**   | X25519 (ECDH) | Curve25519 elliptic curve             |
| **Encryption**     | ChaCha20      | 256-bit keys, 96-bit nonces           |
| **Authentication** | Poly1305      | 128-bit tags                          |
| **Random**         | crypto/rand   | OS-provided CSPRNG                    |

#### Key Hierarchy

```
User Passcode (entered by user)
       │
       ▼
  [Argon2id]
       │
       ▼
Session Key (32 bytes)
       │
       ├──────────────────┐
       │                  │
       ▼                  ▼
 Handshake Key      Handshake Key
  (Initiator)        (Responder)
       │                  │
       └────────┬─────────┘
                │
                ▼
          [Noise DH]
                │
                ├──────────────┬────────────────┐
                │              │                │
                ▼              ▼                ▼
          Send Key       Receive Key     Chaining Key
         (ChaCha20)      (ChaCha20)        (for rekeying)
```

### Security Properties

#### Confidentiality

- All file data encrypted with ChaCha20 (256-bit keys)
- Relay server cannot decrypt traffic (no keys)
- Passive observers see only encrypted bytes

#### Integrity

- Poly1305 MAC on every frame (128-bit tag)
- Tampering detected immediately
- Connection terminated on integrity failure

#### Authentication

- Mutual authentication during handshake
- Both parties prove knowledge of passcode
- Prevents impersonation attacks

#### Forward Secrecy

- Ephemeral X25519 keypairs per session
- Session keys destroyed after handshake
- Past sessions remain secure even if passcode leaked

#### Replay Protection

- Counter-based nonces
- Each frame has unique nonce
- Out-of-order or duplicate frames rejected

## Performance Characteristics

### Latency Profile

```
Operation              Latency    Notes
────────────────────────────────────────────────────────
Session Creation       1-5 ms     Depends on network RTT
Key Derivation         ~100 ms    Intentional (Argon2id)
Noise Handshake        50-200 ms  2 RTT + crypto
File List (100 files)  10-50 ms   1 RTT + encryption
File Read (1 MB)       50-500 ms  Depends on bandwidth
```

### Throughput

ChaCha20 encryption is fast (1-4 GB/s on modern CPUs), so throughput is primarily limited by:

1. **Network bandwidth**: The bottleneck in most scenarios
2. **WebSocket overhead**: Minimal (~2-5% overhead)
3. **Frame size**: 1 MB maximum per frame

### Memory Usage

```
Component           Base Memory    Per Connection
────────────────────────────────────────────────────
Relay Server        ~5 MB          ~100 KB
Sharer (Idle)       ~10 MB         N/A
Receiver (Idle)     ~10 MB         N/A
TUI Browser         +5 MB          N/A
File Buffers        Variable       Up to 10 MB
```

### Scalability

**Relay Server**: Can handle thousands of concurrent sessions on modest hardware (2 CPU cores, 1 GB RAM).

**Limitations**:

- One receiver per session (current implementation)
- Sessions stored in memory (not persistent)
- No horizontal scaling (single instance)

## Design Decisions & Trade-offs

### Why WebSocket for Relay?

**Choice**: WebSocket over raw TCP

**Rationale**:

- Works through HTTP proxies and firewalls
- Automatic keep-alive and reconnection
- Wide support across platforms
- Minimal overhead for binary frames

**Trade-off**: Slightly higher latency vs raw TCP, but better NAT traversal

### Why Noise Protocol?

**Choice**: Noise XX pattern over TLS

**Rationale**:

- Simpler implementation (~500 lines vs thousands)
- Custom authentication flow (passcode-based)
- No certificate infrastructure needed
- Perfect forward secrecy built-in

**Trade-off**: Custom protocol requires more careful auditing

### Why Argon2id?

**Choice**: Argon2id over PBKDF2 or bcrypt

**Rationale**:

- Memory-hard (resistant to GPU/ASIC attacks)
- Hybrid design (resistant to both side-channel and time-memory trade-off attacks)
- Tunable parameters for future security increases

**Trade-off**: Slower key derivation (~100ms) vs instant, but intentional for security

### Why No FUSE/Filesystem Mounting?

**Choice**: TUI browser over FUSE mounting

**Rationale**:

- Cross-platform (Windows support difficult with FUSE)
- Simpler security model (explicit actions)
- No kernel module dependencies
- Easier to audit and maintain

**Trade-off**: Less transparent to user, but more explicit security

### Why Session-Based (No Accounts)?

**Choice**: Ephemeral sessions over persistent accounts

**Rationale**:

- No user database to secure
- No password reset flows
- Natural expiration (24 hours)
- Minimal metadata retention

**Trade-off**: Less convenient for repeated use, but better privacy

### Why Go?

**Choice**: Go over Rust, C++, Python

**Rationale**:

- Memory-safe (no manual memory management)
- Strong standard library (crypto, networking)
- Easy cross-compilation
- Good performance for I/O-bound workloads
- Fast compilation

**Trade-off**: Larger binary size vs Rust/C++, but acceptable for application

## Extending Orb

### Adding New Protocol Operations

To add a new filesystem operation:

1. **Define Operation Code** in `pkg/protocol/protocol.go`:

   ```go
   const (
       OpCopy = 0x08  // New operation
   )
   ```

2. **Implement Handler** in `internal/filesystem/secure_fs.go`:

   ```go
   func (fs *SecureFS) Copy(src, dst string) error {
       // Validate both paths
       // Perform operation
       // Return result
   }
   ```

3. **Add Protocol Marshaling** in `pkg/protocol/protocol.go`:

   ```go
   type CopyRequest struct {
       Source      string
       Destination string
   }
   ```

4. **Handle in Tunnel** in `internal/tunnel/tunnel.go`:
   ```go
   case OpCopy:
       req := parseCopyRequest(frame)
       err := filesystem.Copy(req.Source, req.Destination)
       sendResponse(err)
   ```

### Adding Authentication Methods

Current implementation uses passcode-only. To add alternative auth:

1. **Extend Session** in `internal/session/session.go`:

   ```go
   type AuthMethod interface {
       Verify(credentials []byte) bool
   }
   ```

2. **Implement Method**:

   ```go
   type PKIAuth struct {
       publicKey []byte
   }

   func (p *PKIAuth) Verify(signature []byte) bool {
       // Verify signature
   }
   ```

3. **Integrate with Handshake** in `internal/crypto/noise.go`

### Adding Compression

To add optional compression:

1. **Extend Protocol** with compression flag
2. **Compress Before Encryption**:
   ```go
   compressed := zstd.Compress(plaintext)
   encrypted := encrypt(compressed)
   ```
3. **Decompress After Decryption** on receiver

**Note**: Compression before encryption can leak information about plaintext (CRIME/BREACH attacks). Only use with careful consideration.

## Operational Considerations

### Deployment Patterns

#### Pattern 1: Public Relay

```
Internet
    │
    ├─── [Relay Server (VPS)]
    │         │
    │         ├──── Sharer A (via WebSocket)
    │         ├──── Sharer B (via WebSocket)
    │         └──── Receivers (via WebSocket)
    │
```

**Pros**: Single relay, easy for users
**Cons**: Relay is central point of failure, bandwidth costs

#### Pattern 2: Self-Hosted Relay

```
Sharer's Network
    │
    ├─── [Relay Server (local)]
    │         │
    │         └──── Receiver (via Internet)
```

**Pros**: Full control, no external dependencies
**Cons**: Requires port forwarding or VPN

#### Pattern 3: Corporate Deployment

```
Corporate Network
    │
    ├─── [Load Balancer]
    │         │
    │         ├─── Relay Instance 1
    │         ├─── Relay Instance 2
    │         └─── Relay Instance 3
    │
    └─── Internal Sharers/Receivers
```

**Pros**: High availability, enterprise control
**Cons**: Requires infrastructure, session stickiness needed

### Monitoring & Observability

#### Key Metrics to Monitor

**Relay Server**:

- Active WebSocket connections
- Session creation rate
- Data throughput (bytes/sec)
- Connection errors
- Average session duration

**Application**:

- Handshake success rate
- Authentication failures (potential attacks)
- File operation latencies
- Encryption/decryption throughput

#### Logging Best Practices

**What to Log**:

- Connection events (connect, disconnect)
- Session creation/expiration
- Authentication failures (for rate limiting)
- Error conditions

**What NOT to Log**:

- Passcodes (security risk)
- File names (privacy concern)
- File contents (obviously)
- Session IDs in plaintext (if logs are public)

### Disaster Recovery

#### Session State Loss

**Problem**: Relay server crashes, all sessions lost

**Mitigation**: Sessions are ephemeral by design, users simply recreate

#### Relay Compromise

**Problem**: Attacker gains relay server access

**Impact**: Attacker can:

- See connection metadata (IPs, timing)
- Drop connections (DoS)
- Cannot decrypt traffic (zero-knowledge)

**Recovery**: Rotate to new relay server, alert users

## Testing Strategy

### Unit Tests

Focus areas:

- Cryptographic primitives (key derivation, encryption)
- Path sanitization (boundary checks)
- Protocol parsing (malformed frames)
- Session management (expiration, locking)

### Integration Tests

Test scenarios:

- Complete handshake flow
- File operations over encrypted tunnel
- Authentication failures and rate limiting
- Session expiration

### Security Tests

Critical tests:

- Path traversal attempts
- Symlink escape attempts
- Replay attack simulation
- Malformed protocol frames
- Timing attack resistance (constant-time ops)

### Performance Tests

Benchmarks:

- Key derivation time (should be ~100ms)
- Encryption throughput (MB/s)
- File transfer speed
- Concurrent session handling (relay)

## Future Enhancements

### Potential Improvements

1. **Resume Support**: Allow interrupted file transfers to resume
2. **Multi-Receiver**: Enable multiple receivers per session
3. **Persistent Sessions**: Optional session persistence across relay restarts
4. **IPv6 Support**: Full native IPv6 support
5. **Compression**: Optional transparent compression (with security considerations)
6. **Mobile Clients**: iOS/Android native apps
7. **Browser Client**: WebAssembly-based browser client

### Architectural Constraints

Any enhancement must maintain:

- Zero-knowledge relay (no plaintext access)
- End-to-end encryption
- No persistent user state
- Minimal attack surface

## References

### Cryptography

- [Noise Protocol Framework](https://noiseprotocol.org/)
- [RFC 8439: ChaCha20-Poly1305](https://www.rfc-editor.org/rfc/rfc8439)
- [RFC 9106: Argon2](https://www.rfc-editor.org/rfc/rfc9106)
- [Curve25519](https://cr.yp.to/ecdh.html)

### Implementation

- [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [Bubble Tea TUI](https://github.com/charmbracelet/bubbletea)

### Security Research

- [Cryptographic Right Answers](https://latacora.micro.blog/2018/04/03/cryptographic-right-answers.html)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Password Hashing Competition](https://www.password-hashing.net/)

---

**Document Version**: 1.0  
**Last Updated**: January 15, 2026  
**Status**: Living Document

For security issues, please see [security policy](security.md).
