# Orb Security Architecture

## Threat Model

Orb is designed to be secure against the following threats:

### Network Attackers

- **Passive Eavesdropping**: All data is encrypted end-to-end
- **Active Man-in-the-Middle**: Noise Protocol provides mutual authentication
- **Replay Attacks**: Unique nonces prevent replay
- **Session Hijacking**: Authentication required for every session

### Malicious Relay Server

- **Data Interception**: Relay only sees encrypted bytes
- **Authentication Bypass**: Authentication happens peer-to-peer
- **Metadata Analysis**: Minimal metadata exposed

### Malicious Clients

- **Path Traversal**: All paths are sanitized and validated
- **Symlink Exploitation**: Symlinks are resolved and checked
- **Resource Exhaustion**: Rate limiting and size limits enforced
- **Privilege Escalation**: Sandboxed filesystem operations

### Brute Force Attacks

- **Passcode Guessing**: Argon2id makes each attempt expensive (>100ms)
- **Session Enumeration**: Generic error messages prevent information leakage
- **Account Lockout**: Sessions lock after 5 failed attempts

## Cryptographic Primitives

### Key Derivation: Argon2id

**Purpose**: Derive encryption keys from user passcode

**Parameters**:

- Time cost: 3 iterations
- Memory cost: 64 MB
- Parallelism: 4 threads
- Output: 32 bytes

**Security Properties**:

- Memory-hard: Resistant to GPU/ASIC attacks
- Side-channel resistant: Constant-time operations
- Tunable: Parameters can be adjusted for stronger security

**Implementation**:

```go
key := argon2.IDKey(passcode, sessionID, 3, 64*1024, 4, 32)
```

### Key Exchange: Noise Protocol (X25519)

**Purpose**: Establish secure channel with mutual authentication

**Pattern**: Simplified Noise_XX with pre-shared key

**Steps**:

1. Both parties generate ephemeral X25519 key pairs
2. Initiator sends ephemeral public key + encrypted auth
3. Responder validates auth, sends ephemeral public key + encrypted auth
4. Both derive shared transport keys

**Security Properties**:

- Perfect forward secrecy: Ephemeral keys ensure past sessions remain secure
- Mutual authentication: Both parties prove knowledge of passcode
- Identity hiding: No static identities revealed
- Post-compromise security: Future sessions remain secure

**Implementation**: See `internal/crypto/noise.go`

### Transport Encryption: ChaCha20-Poly1305

**Purpose**: Authenticated encryption of tunnel traffic

**Variant**: XChaCha20-Poly1305 (extended nonce)

**Security Properties**:

- Authenticated encryption: Confidentiality + integrity
- Large nonce space: 192-bit nonces prevent reuse
- Fast in software: No need for hardware acceleration
- Constant-time: Resistant to timing attacks

**Frame Format**:

```
[24-byte nonce][ciphertext][16-byte auth tag]
```

**Implementation**:

```go
cipher, _ := chacha20poly1305.NewX(key)
ciphertext := cipher.Seal(nonce, nonce, plaintext, nil)
```

## Attack Surface Analysis

### Client Attack Surface

**Input Validation**:

- Session IDs: Validated format and length
- Passcodes: No injection possible (used as bytes)
- File paths: Sanitized and bounded to shared directory
- Frame types: Validated against whitelist
- Frame sizes: Limited to 1 MB maximum

**State Management**:

- Session expiration: 24-hour timeout
- Connection timeout: 60-second read timeout
- Failed attempt tracking: Lock after 5 failures
- Replay protection: Nonce counter increments

**Resource Limits**:

- Message size: 2 MB WebSocket limit
- File read: 10 MB per read operation
- Memory: Bounded buffer sizes
- Connections: Per-session limits enforced

### Relay Server Attack Surface

**Minimal Trusted Computing Base**:

- No authentication logic
- No decryption capability
- No file storage
- No logging of payloads

**DoS Protection**:

- Connection timeout: 60 seconds idle
- Session cleanup: Automatic expiration
- Rate limiting: Per-IP connection limits (TODO)
- Resource limits: WebSocket buffer sizes

**Information Leakage**:

- Generic error messages
- No session enumeration
- No timing attacks (constant-time comparison)
- No metadata logging

## Security Checklist

### Cryptography 

- [x] Key derivation uses Argon2id
- [x] Handshake uses Noise Protocol
- [x] Transport uses ChaCha20-Poly1305
- [x] Random numbers use crypto/rand
- [x] Keys are zeroized after use
- [x] Constant-time comparisons

### Authentication 

- [x] Mutual authentication required
- [x] No static secrets
- [x] No passwords sent in plaintext
- [x] Failed attempt rate limiting
- [x] Session lockout after failures

### Authorization 

- [x] Path traversal prevention
- [x] Symlink escape detection
- [x] Read-only mode support
- [x] No remote execution
- [x] Sandboxed operations

### Network Security 

- [x] All data encrypted end-to-end
- [x] Replay protection
- [x] Nonce uniqueness
- [x] Connection timeouts
- [x] NAT traversal

### Privacy 

- [x] Relay server is blind
- [x] No passcode logging
- [x] No filename logging
- [x] No content logging
- [x] Minimal metadata

## Known Limitations

### Not Protected Against

1. **Endpoint Compromise**: If attacker has access to client machine, they can access shared files
2. **Timing Analysis**: Large-scale correlation attacks may reveal session patterns
3. **DoS on Relay**: Public relay can be overwhelmed (deploy rate limiting)
4. **Malicious Files**: Orb does not scan for malware in shared files
5. **Social Engineering**: If user shares passcode with attacker, security is compromised

### Assumptions

1. **Client Security**: Users' machines are not compromised
2. **Relay Availability**: Relay server is available and honest-but-curious
3. **Time Synchronization**: Clocks are reasonably synchronized for timeouts
4. **Memory Safety**: Go runtime provides memory safety
5. **Crypto Libraries**: golang.org/x/crypto is correct and secure

## Security Updates

### Versioning

Orb follows semantic versioning:

- **Major**: Breaking changes, including security architecture changes
- **Minor**: New features, backward compatible
- **Patch**: Bug fixes, security fixes

### Update Policy

Security updates are released immediately and backported to supported versions.

### Supported Versions

| Version | Supported           |
| ------- | ------------------- |
| 1.x     | Yes              |
| < 1.0   | No (pre-release) |

## Reporting Vulnerabilities

**DO NOT** create public GitHub issues for security vulnerabilities.

**Email**: security@orb.example.com

**PGP Key**: [Include PGP key]

**Response Time**:

- Acknowledgment: Within 48 hours
- Assessment: Within 7 days
- Fix: Within 30 days (depending on severity)

**Disclosure Policy**:

- We follow coordinated disclosure
- Vulnerabilities are disclosed 90 days after fix or public announcement
- Critical vulnerabilities may be disclosed sooner

## Security Audit Status

**Status**: Not yet audited

**Recommended**: Before using in production, conduct a security audit of:

- Cryptographic implementation
- Network protocol
- Filesystem sandboxing
- Attack surface

**Bounty Program**: Not yet available

## References

- [Noise Protocol Framework](https://noiseprotocol.org/)
- [RFC 8439: ChaCha20-Poly1305](https://www.rfc-editor.org/rfc/rfc8439)
- [RFC 9106: Argon2](https://www.rfc-editor.org/rfc/rfc9106)
- [OWASP Secure Coding Practices](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)

---

Last Updated: 2026-01-14
