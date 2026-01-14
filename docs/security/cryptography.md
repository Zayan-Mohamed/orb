# Cryptography

Detailed technical documentation of Orb's cryptographic implementation.

## Overview

Orb uses a defense-in-depth cryptographic architecture:

1. **Key Derivation**: Argon2id derives keys from passcode
2. **Key Exchange**: Noise Protocol establishes secure channel
3. **Transport Encryption**: ChaCha20-Poly1305 encrypts all data

All cryptographic operations use vetted implementations from `golang.org/x/crypto`.

## Key Derivation: Argon2id

### Purpose

Convert user-memorable passcode into cryptographic key material:

- Input: Short passcode (random string)
- Output: 32-byte cryptographic key
- Properties: Memory-hard, GPU-resistant

### Parameters

```go
argon2.IDKey(
    password,     // Session passcode
    salt,         // Session ID (acts as salt)
    time: 3,      // Iterations
    memory: 64*1024, // 64 MB memory
    threads: 4,   // Parallel threads
    keyLen: 32,   // 256-bit key output
)
```

### Security Properties

**Memory-Hardness:**

- Requires 64 MB of RAM per attempt
- Prevents GPU/ASIC attacks
- Makes brute force expensive

**Tunable Parameters:**

- Can increase for stronger security
- Current settings: ~100ms per derivation
- Balances security vs usability

**Side-Channel Resistance:**

- Constant-time operations
- No data-dependent branches
- Resistant to timing attacks

### Implementation Details

```go
// internal/crypto/crypto.go
func DeriveKey(passcode, sessionID []byte) []byte {
    return argon2.IDKey(
        passcode,
        sessionID,
        3,      // time
        64*1024, // memory in KB
        4,      // threads
        32,     // key length
    )
}
```

The derived key is used as input to the Noise Protocol handshake.

## Noise Protocol

### Overview

Noise Protocol provides:

- Secure key exchange
- Mutual authentication
- Perfect forward secrecy
- Identity hiding

### Pattern

Orb uses a simplified Noise_XX pattern with pre-shared key (passcode):

```
-> e
<- e, ee, s, es
-> s, se
```

Simplified for Orb:

1. Initiator sends ephemeral public key + encrypted challenge
2. Responder validates, sends ephemeral public key + encrypted response
3. Both derive transport keys from shared secret

### Handshake Flow

**Initialization:**

```go
// Generate ephemeral X25519 key pair
privateKey, publicKey := generateKeyPair()

// Derive base key from passcode
baseKey := argon2.IDKey(passcode, sessionID, 3, 64*1024, 4, 32)
```

**Message 1 (Initiator → Responder):**

```go
// Initiator generates ephemeral key
ephemeralPriv, ephemeralPub := generateKeyPair()

// Encrypt authentication challenge
challenge := "orb-handshake-v1"
encrypted := encrypt(baseKey, challenge)

// Send: ephemeral public key + encrypted challenge
message1 := ephemeralPub || encrypted
```

**Message 2 (Responder → Initiator):**

```go
// Responder receives ephemeral pub
initiatorEphemeral := message1[:32]

// Decrypt and validate challenge
decrypted := decrypt(baseKey, encrypted)
if decrypted != "orb-handshake-v1" {
    return error("authentication failed")
}

// Generate own ephemeral key
ephemeralPriv, ephemeralPub := generateKeyPair()

// Compute shared secret
sharedSecret := x25519(ephemeralPriv, initiatorEphemeral)

// Encrypt response
response := "orb-handshake-v1-ok"
encrypted := encrypt(baseKey, response)

// Send: ephemeral public key + encrypted response
message2 := ephemeralPub || encrypted
```

**Key Derivation:**

```go
// Both parties compute shared secret
sharedSecret := x25519(myPrivate, theirPublic)

// Derive transport keys
initiatorToResponder := hkdf(sharedSecret, "initiator_to_responder")
responderToInitiator := hkdf(sharedSecret, "responder_to_initiator")

// Initiator uses: send=i2r, recv=r2i
// Responder uses: send=r2i, recv=i2r
```

### Security Properties

**Perfect Forward Secrecy:**

- Ephemeral keys discarded after handshake
- Compromise of passcode doesn't reveal past sessions
- Each session has unique keys

**Mutual Authentication:**

- Both parties prove knowledge of passcode
- Prevents unauthorized connections
- Binds session to passcode

**Identity Hiding:**

- No static public keys exchanged
- Relay cannot identify parties
- Passcode never sent in plaintext

**Post-Compromise Security:**

- Future sessions secure even if one session compromised
- New ephemeral keys each session

### Implementation

See [internal/crypto/noise.go](../development/api.md#noise-protocol) for full implementation.

## Transport Encryption: ChaCha20-Poly1305

### Purpose

After handshake, all tunnel traffic is encrypted using authenticated encryption:

- Confidentiality: ChaCha20 stream cipher
- Authenticity: Poly1305 MAC
- Combined: AEAD construction

### Variant

**XChaCha20-Poly1305:**

- Extended nonce: 192 bits (vs 96 bits in ChaCha20)
- Allows random nonce generation
- No nonce reuse risk with reasonable message counts

### Frame Structure

Each frame is encrypted:

```
+-------------+------------------+-----------+
| Nonce (24B) | Ciphertext (N B) | Tag (16B) |
+-------------+------------------+-----------+
```

**Components:**

- **Nonce**: 24-byte unique value per message
- **Ciphertext**: Encrypted plaintext
- **Tag**: 16-byte authentication tag

### Nonce Generation

**Strategy: Counter-based**

```go
type AEAD struct {
    cipher   cipher.AEAD
    sendNonce uint64
    recvNonce uint64
}

func (a *AEAD) Encrypt(plaintext []byte) []byte {
    // Increment counter
    nonce := a.sendNonce
    a.sendNonce++

    // Convert to 24-byte nonce
    var nonceBytes [24]byte
    binary.BigEndian.PutUint64(nonceBytes[16:], nonce)

    // Encrypt
    return cipher.Seal(nonceBytes[:], nonceBytes[:], plaintext, nil)
}
```

**Properties:**

- Unique per message
- Never repeats (64-bit counter)
- No randomness needed
- Synchronized between parties

### Encryption Operation

```go
func Encrypt(key, plaintext []byte) []byte {
    // Create cipher
    cipher, _ := chacha20poly1305.NewX(key)

    // Generate nonce
    nonce := generateNonce()

    // Encrypt: nonce || ciphertext || tag
    ciphertext := cipher.Seal(nonce, nonce, plaintext, nil)

    return ciphertext
}
```

### Decryption Operation

```go
func Decrypt(key, ciphertext []byte) ([]byte, error) {
    // Create cipher
    cipher, _ := chacha20poly1305.NewX(key)

    // Extract nonce
    nonce := ciphertext[:24]

    // Decrypt and authenticate
    plaintext, err := cipher.Open(nil, nonce, ciphertext[24:], nil)
    if err != nil {
        return nil, errors.New("decryption failed")
    }

    return plaintext, nil
}
```

### Security Properties

**Authenticated Encryption:**

- Ciphertext integrity guaranteed
- Detects any tampering
- Prevents bit-flipping attacks

**Nonce Uniqueness:**

- 64-bit counter ensures no reuse
- Can encrypt 2^64 messages safely
- No birthday bound concerns

**Constant-Time:**

- No data-dependent branches
- Resistant to timing attacks
- Safe against side-channel analysis

**Performance:**

- Fast in software (no AES-NI needed)
- ~5 cycles/byte on modern CPUs
- Suitable for high-throughput

## Key Management

### Key Lifecycle

1. **Generation**: Derived from passcode + session ID
2. **Usage**: Handshake and transport encryption
3. **Rotation**: Not supported (create new session)
4. **Deletion**: Zeroized after use

### Key Storage

**In Memory Only:**

- Keys never written to disk
- Held only for session duration
- Cleared when connection closes

**Zeroization:**

```go
// Clear sensitive data
defer func() {
    for i := range key {
        key[i] = 0
    }
}()
```

### Key Separation

Different keys for different purposes:

- Handshake: baseKey from Argon2id
- Transport: Derived from Noise handshake
- Send/Receive: Separate keys for each direction

## Cryptographic Dependencies

### Go Crypto Library

All primitives from `golang.org/x/crypto`:

```go
import (
    "golang.org/x/crypto/argon2"
    "golang.org/x/crypto/chacha20poly1305"
    "golang.org/x/crypto/curve25519"
)
```

**Why golang.org/x/crypto?**

- Maintained by Go team
- Constant-time implementations
- Well-audited
- Cross-platform

### Random Number Generation

```go
import "crypto/rand"

// Generate random bytes
func Random(n int) []byte {
    b := make([]byte, n)
    rand.Read(b)
    return b
}
```

**Properties:**

- Cryptographically secure
- Uses OS entropy source (/dev/urandom, CryptGenRandom)
- Suitable for keys, nonces, session IDs

## Security Considerations

### What Orb Protects Against

✅ **Eavesdropping**: All data encrypted
✅ **Tampering**: Authentication tags prevent modification
✅ **Replay**: Nonce prevents replay attacks
✅ **MITM**: Mutual authentication via passcode
✅ **Brute Force**: Argon2id makes attempts expensive

### What Orb Does NOT Protect Against

❌ **Endpoint Compromise**: If client hacked, files accessible
❌ **Weak Passcodes**: User must choose strong passcode
❌ **Social Engineering**: If passcode shared with attacker
❌ **Side Channels**: Advanced attacks (cache timing, etc.)
❌ **Quantum Computers**: X25519 vulnerable to quantum attacks

## Future Enhancements

### Post-Quantum Cryptography

Future versions may add:

- Kyber for key exchange
- Dilithium for signatures
- Hybrid classical + post-quantum

### Hardware Security Modules

Support for:

- TPM-backed key storage
- HSM for key derivation
- Secure enclaves (SGX, etc.)

### Key Rotation

Planned features:

- Periodic re-keying
- Forward-secure ratcheting
- Session resumption

## Testing

### Test Vectors

See [internal/crypto/crypto_test.go](../development/building.md#running-tests) for:

- Argon2id test vectors
- Noise handshake test cases
- Encryption/decryption tests

### Fuzzing

Cryptographic functions are fuzzed:

```bash
go test -fuzz=FuzzEncryption
go test -fuzz=FuzzHandshake
```

## References

- [Noise Protocol Framework](https://noiseprotocol.org/)
- [RFC 8439: ChaCha20-Poly1305](https://www.rfc-editor.org/rfc/rfc8439)
- [RFC 9106: Argon2](https://www.rfc-editor.org/rfc/rfc9106)
- [NaCl: Networking and Cryptography library](https://nacl.cr.yp.to/)
- [Go Cryptography](https://pkg.go.dev/golang.org/x/crypto)

## Next Steps

- Read [Security Best Practices](best-practices.md)
- Understand [Threat Model](threat-model.md)
- Review [Security Overview](security.md)
