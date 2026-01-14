# Threat Model

Understanding Orb's security design and threat analysis.

## Overview

Orb is designed to securely share files between two parties through an untrusted relay server. This document analyzes the threat model and security guarantees.

## Trust Assumptions

### Trusted Components

**✅ Client Endpoints:**

- User's local machines are not compromised
- Operating system is trusted
- Go runtime is memory-safe
- Crypto libraries are correct

**✅ Users:**

- Users protect their passcodes
- Users authenticate each other out-of-band
- Users follow security best practices

### Untrusted Components

**❌ Relay Server:**

- Considered "honest-but-curious"
- May log metadata
- May be compromised
- Cannot break encryption

**❌ Network:**

- Passive eavesdroppers present
- Active attackers may inject/modify
- No assumptions about network security

## Threat Actors

### 1. Network Eavesdropper

**Capabilities:**

- Observe all network traffic
- Record encrypted sessions
- Passive monitoring only

**Goals:**

- Decrypt file contents
- Identify communicating parties
- Map file access patterns

**Protections:**

✅ **End-to-End Encryption:**

```
All data encrypted with ChaCha20-Poly1305
Relay sees only ciphertext
```

✅ **Forward Secrecy:**

```
Ephemeral keys per session
Past sessions secure even if passcode compromised later
```

✅ **Minimal Metadata:**

```
Session IDs are random
No user identifiers transmitted
```

**Remaining Risks:**

- Traffic analysis (connection timing, sizes)
- Correlation attacks over time

### 2. Man-in-the-Middle Attacker

**Capabilities:**

- Intercept network traffic
- Modify messages in transit
- Impersonate endpoints
- Block communication

**Goals:**

- Decrypt communications
- Inject malicious files
- Impersonate sharer or connector

**Protections:**

✅ **Mutual Authentication:**

```
Noise Protocol handshake
Both parties prove knowledge of passcode
Cannot be impersonated without passcode
```

✅ **Authenticated Encryption:**

```
Poly1305 MAC prevents tampering
Modified messages rejected
```

✅ **Replay Protection:**

```
Nonce counter prevents replays
Each message is unique
```

**Remaining Risks:**

- If attacker obtains passcode (social engineering)
- Denial of service (blocking connections)

### 3. Malicious Relay Server

**Capabilities:**

- See all encrypted traffic
- Log session IDs
- Observe connection patterns
- Terminate connections
- Delay messages

**Goals:**

- Decrypt file contents
- Identify users
- Map social connections
- Disrupt service

**Protections:**

✅ **Blind Relay Design:**

```
Relay has no decryption keys
Cannot see plaintext data
End-to-end encryption bypasses relay
```

✅ **No Authentication at Relay:**

```
Relay doesn't know user identities
Session IDs are random
No persistent accounts
```

✅ **Tamper Detection:**

```
Authentication tags detect modification
Encrypted handshake prevents impersonation
```

**Remaining Risks:**

- Metadata analysis (timing, size, patterns)
- Denial of service
- Traffic correlation attacks
- Logging session IDs for enumeration

### 4. Malicious Client

**Capabilities:**

- Send crafted file requests
- Attempt path traversal
- Symlink exploitation
- Resource exhaustion
- Protocol abuse

**Goals:**

- Read files outside shared directory
- Escalate privileges
- Crash sharer
- Consume resources

**Protections:**

✅ **Path Sanitization:**

```go
// internal/filesystem/secure_fs.go
func sanitizePath(base, requested string) (string, error) {
    // Resolve to absolute path
    absRequested := filepath.Join(base, filepath.Clean(requested))

    // Ensure within base directory
    if !strings.HasPrefix(absRequested, base) {
        return "", errors.New("path traversal detected")
    }

    return absRequested, nil
}
```

✅ **Symlink Protection:**

```
Symlinks resolved and validated
Cannot escape shared directory
```

✅ **Resource Limits:**

```
Max file read: 10 MB per request
Max message size: 2 MB
Connection timeouts: 60 seconds
```

✅ **Read-Only Access:**

```
No write operations supported
No file deletion
No remote execution
```

**Remaining Risks:**

- Authorized files may contain malware
- Resource consumption within limits
- Timing side channels

### 5. Brute Force Attacker

**Capabilities:**

- Guess session IDs
- Try multiple passcodes
- Automate connection attempts
- Use GPU/ASIC for cracking

**Goals:**

- Guess correct passcode
- Enumerate valid sessions
- Access unauthorized files

**Protections:**

✅ **Argon2id Key Derivation:**

```
Time: 3 iterations (~100ms per attempt)
Memory: 64 MB per attempt
GPU/ASIC resistant
```

✅ **Rate Limiting:**

```
5 failed attempts = session lockout
Generic error messages
No timing attacks
```

✅ **Strong Randomness:**

```
Session IDs: 12 random characters (62^12 ≈ 3×10^21)
Passcodes: High entropy random strings
```

**Remaining Risks:**

- Weak user-chosen passcodes (if supported in future)
- Distributed attacks across many sessions
- Offline cracking if ciphertext captured (mitigated by Argon2id)

## Attack Scenarios

### Scenario 1: Passive Surveillance

**Attacker:** Intelligence agency, ISP
**Goal:** Identify who is sharing what with whom

**Attack:**

1. Monitor all network traffic
2. Record encrypted sessions
3. Correlate connection patterns
4. Perform traffic analysis

**Orb Protection:**

- ✅ Traffic encrypted end-to-end
- ✅ No plaintext metadata
- ⚠️ Connection timing visible
- ⚠️ Session IDs in WebSocket headers

**Mitigation:**

- Use TLS relay (wss://)
- Mix traffic with cover traffic
- Use Tor for anonymity (advanced)

### Scenario 2: Malicious Relay Operator

**Attacker:** Untrusted relay server operator
**Goal:** Decrypt shared files

**Attack:**

1. Set up malicious relay
2. Log all traffic
3. Attempt to decrypt
4. Break encryption

**Orb Protection:**

- ✅ End-to-end encryption
- ✅ Relay has no keys
- ✅ Perfect forward secrecy
- ✅ Cannot decrypt without passcode

**Mitigation:**

- Self-host relay
- Use trusted relay operators
- Regularly rotate sessions

### Scenario 3: Compromised Endpoint

**Attacker:** Malware on client machine
**Goal:** Steal files or passcodes

**Attack:**

1. Keylogger captures passcode
2. Screen capture reveals session ID
3. Memory dump extracts keys
4. Direct file access

**Orb Protection:**

- ❌ Cannot protect against endpoint compromise
- ❌ If machine is pwned, game over

**Mitigation:**

- Keep systems patched
- Use antivirus/EDR
- Encrypt disk
- Physical security

### Scenario 4: Social Engineering

**Attacker:** Phisher or impersonator
**Goal:** Trick user into sharing passcode

**Attack:**

1. Impersonate legitimate user
2. Request session credentials
3. Use credentials to access files

**Orb Protection:**

- ✅ Technical controls work correctly
- ❌ Cannot prevent user from sharing passcode

**Mitigation:**

- Verify identity out-of-band
- Use encrypted channels
- Never share passcodes publicly
- Security awareness training

### Scenario 5: Insider Threat

**Attacker:** Authorized user with malicious intent
**Goal:** Exfiltrate sensitive files

**Attack:**

1. User has legitimate access
2. Shares sensitive directory
3. Downloads to unauthorized device
4. Leaks or sells data

**Orb Protection:**

- ❌ Cannot prevent authorized access abuse
- ✅ Can log and audit activity

**Mitigation:**

- Principle of least privilege
- Audit sharing logs
- Monitor for anomalies
- Data loss prevention tools

## Security Boundaries

### What Orb Protects

| Threat                | Protection                  |
| --------------------- | --------------------------- |
| Network eavesdropping | ✅ End-to-end encryption    |
| Man-in-the-middle     | ✅ Mutual authentication    |
| Malicious relay       | ✅ Zero-knowledge design    |
| Path traversal        | ✅ Sandboxing               |
| Replay attacks        | ✅ Nonce counters           |
| Brute force           | ✅ Argon2id + rate limiting |
| Data tampering        | ✅ Authentication tags      |
| Session hijacking     | ✅ Cryptographic binding    |

### What Orb Does NOT Protect

| Threat              | Reason                |
| ------------------- | --------------------- |
| Endpoint compromise | Out of scope          |
| Weak passcodes      | User responsibility   |
| Social engineering  | User error            |
| Malware in files    | Content-agnostic      |
| Quantum computers   | Classical crypto only |
| Traffic analysis    | Metadata leakage      |
| Insider threats     | Authorized access     |
| Physical access     | Out of scope          |

## Cryptographic Security

### Security Levels

**Key Sizes:**

- Argon2id: 256-bit output
- X25519: 128-bit security
- ChaCha20: 256-bit keys
- Poly1305: 128-bit security

**Overall Security Level:** ~128 bits

**Equivalent to:**

- RSA 3072-bit
- AES-128
- Sufficient for SECRET classification (NSA Suite B)

### Known Weaknesses

**Not Quantum-Resistant:**

- X25519 broken by Shor's algorithm
- ChaCha20 weakened by Grover's algorithm
- Post-quantum migration needed

**Side-Channel Leakage:**

- Timing attacks mitigated but not eliminated
- Cache timing possible on shared CPUs
- Power analysis on embedded devices

**Implementation Bugs:**

- Orb not formally verified
- Potential for subtle bugs
- Regular security audits recommended

## Compliance Considerations

### Data Protection Regulations

**GDPR:**

- ✅ Encryption in transit
- ✅ Minimal data collection
- ✅ Right to erasure (stop sharing)
- ✅ Data minimization

**HIPAA:**

- ✅ Technical safeguards
- ⚠️ Organizational policies required
- ⚠️ Business associate agreements needed

**CCPA:**

- ✅ Consumer privacy
- ✅ Data security
- ✅ No sale of data

### Export Controls

**Encryption Strength:**

- ChaCha20: 256-bit keys
- May be subject to export controls
- Check local regulations

## Risk Assessment

### High Risk

❌ **Endpoint compromise** → Complete system compromise
❌ **Passcode phishing** → Unauthorized access
❌ **Physical device theft** → Data exposure if disk not encrypted

### Medium Risk

⚠️ **Traffic analysis** → Metadata leakage
⚠️ **Malicious relay** → DoS, traffic analysis
⚠️ **Weak passcodes** → Brute force possible

### Low Risk

✅ **Network eavesdropping** → Encrypted
✅ **MITM attacks** → Authenticated
✅ **Path traversal** → Mitigated

## Future Enhancements

### Planned Improvements

1. **Post-Quantum Crypto:**

   - Add Kyber for key exchange
   - Hybrid classical + PQ

2. **Anonymity:**

   - Tor integration
   - Mix networks
   - Traffic padding

3. **Metadata Protection:**

   - Onion routing
   - Traffic obfuscation
   - Timing randomization

4. **Hardware Security:**
   - TPM integration
   - HSM support
   - Secure enclaves

## Conclusion

Orb provides strong security against network attackers and malicious relays through:

- End-to-end encryption
- Mutual authentication
- Zero-knowledge relay design
- Filesystem sandboxing

However, users must:

- Protect passcodes
- Verify peer identity
- Secure endpoints
- Follow best practices

## Next Steps

- Read [Cryptography Details](cryptography.md)
- Review [Best Practices](best-practices.md)
- Check [Security Overview](../security.md)
