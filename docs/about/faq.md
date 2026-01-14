# Frequently Asked Questions

Common questions about Orb.

## General

### What is Orb?

Orb is a zero-trust, end-to-end encrypted file sharing tool that allows you to securely share directories between two parties through an untrusted relay server.

### How is Orb different from other file sharing tools?

- **Zero-trust**: Relay server cannot see your files
- **End-to-end encrypted**: Only you and recipient can decrypt
- **No account required**: Session-based, no registration
- **Open source**: Auditable security
- **Cross-platform**: Works on Linux, macOS, Windows

### Is Orb free?

Yes, Orb is open-source under MIT license.

## Security

### Can the relay server see my files?

No. All files are encrypted end-to-end using ChaCha20-Poly1305. The relay only sees encrypted bytes.

### How secure is Orb?

Orb uses industry-standard cryptography:

- Argon2id for key derivation
- Noise Protocol for handshake
- ChaCha20-Poly1305 for encryption
- X25519 for key exchange

### What if someone steals my passcode?

If an attacker obtains your passcode, they can access that session. However:

- Sessions expire after 24 hours
- Each session has unique credentials
- No persistent access
- You can stop sharing immediately

### Is Orb quantum-resistant?

No. Orb uses classical cryptography (X25519, ChaCha20) which is vulnerable to quantum attacks. Post-quantum cryptography is planned for future versions.

## Usage

### Can multiple people connect to one share?

Currently, one session = one connection. For multiple users, create multiple sharing sessions.

### Can I upload files to a share?

No. Orb is read-only. Only the connector can download files, not upload or modify them.

### How long do sessions last?

Sessions expire after 24 hours. After that, you need to create a new session.

### Can I resume interrupted downloads?

Not currently. Downloads must complete in one session. Resume capability is planned for future versions.

### Does Orb work over the internet?

Yes, as long as both parties can reach the relay server. The relay can be:

- localhost (same machine, for testing)
- Local network
- Public internet server
- Self-hosted server

### What's the maximum file size?

No hard limit, but:

- Per-read limit: 10 MB
- Large files take longer
- No progress bar (yet)
- No resume (yet)

## Technical

### What ports does Orb use?

- Relay server: Default 8080 (WebSocket)
- Clients: Only outbound connections, no ports opened

### Does Orb work behind NAT?

Yes. Clients only need outbound connectivity to the relay server. No port forwarding required.

### Does Orb work behind corporate firewall?

Usually yes, if WebSocket connections are allowed. For restrictive firewalls:

- Use standard ports (80/443)
- Use WSS (encrypted WebSocket)
- Check with IT department

### Can I self-host a relay?

Yes! Just run:

```bash
orb relay --host 0.0.0.0 --port 8080
```

See [Relay Server Deployment](../deployment/relay-server.md) for details.

### What programming language is Orb written in?

Go (Golang) 1.21+

## Troubleshooting

### I get "connection refused" error

The relay server is not reachable. Check:

- Is the relay running?
- Is the URL correct?
- Is there a firewall blocking?
- Try `curl http://relay:8080`

### I get "authentication failed" error

Wrong passcode or session ID. Double-check:

- Session ID is correct
- Passcode is exact (case-sensitive)
- No extra spaces
- Session hasn't expired

### Downloads are very slow

- Check network speed
- Use local/nearby relay
- Large files take time
- No progress bar (yet)

### TUI browser doesn't work

- Use modern terminal
- Enable UTF-8 encoding
- Minimum 80x24 terminal size
- Try different terminal emulator

## Comparison

### Orb vs. SFTP

**Orb:**

- No server setup
- No accounts
- Temporary sessions
- TUI browser

**SFTP:**

- Requires SSH server
- User accounts
- Persistent access
- Any SFTP client

### Orb vs. Magic Wormhole

**Similarities:**

- Encrypted transfer
- Session-based
- No accounts

**Differences:**

- Orb: Directory browsing, TUI
- Wormhole: Single file, CLI only
- Orb: WebSocket relay
- Wormhole: PAKE over Tor

### Orb vs. Cloud Storage

**Orb:**

- No upload to cloud
- Direct peer-to-peer
- No storage fees
- Temporary
- Full control

**Cloud Storage:**

- Files stored centrally
- Persistent storage
- Subscription costs
- Easier sharing
- Provider has access

## Development

### Can I contribute?

Yes! See [Contributing Guide](../development/contributing.md).

### How do I report bugs?

[GitHub Issues](https://github.com/Zayan-Mohamed/orb/issues)

### How do I report security issues?

Email: security@orb.example.com

Do not create public issues for security vulnerabilities.

### What's the roadmap?

See [GitHub Projects](https://github.com/Zayan-Mohamed/orb/projects) and [Changelog](changelog.md).

## Legal

### What license is Orb under?

MIT License. See [License](license.md).

### Can I use Orb commercially?

Yes, MIT license allows commercial use.

### Is Orb GDPR compliant?

Orb provides technical controls (encryption, minimal data collection). Full GDPR compliance requires organizational policies.

### Can I use Orb for HIPAA data?

Orb provides encryption and access controls, but:

- Full HIPAA compliance requires more
- Consult with compliance expert
- Consider additional safeguards

## Still Have Questions?

- [Documentation](../index.md)
- [GitHub Discussions](https://github.com/Zayan-Mohamed/orb/discussions)
- [Security](../security/security.md)
