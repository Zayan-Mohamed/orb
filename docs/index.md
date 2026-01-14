# Welcome to Orb

<div align="center">
  <h1>🛡️ Orb - Zero-Trust Folder Tunneling</h1>
  <p><strong>Secure file sharing with end-to-end encryption</strong></p>
  <p>No accounts • No cloud • No port forwarding</p>
</div>

---

## What is Orb?

Orb is a **production-ready** zero-trust file sharing tool that allows you to securely share a local folder across the internet using end-to-end encryption. It's designed for security professionals, developers, and teams who need to share files without relying on cloud services.

### Key Features

✅ **Zero-Trust Architecture** - The relay server never sees plaintext data  
✅ **Strong Cryptography** - Argon2id, Noise Protocol, ChaCha20-Poly1305  
✅ **Cross-Platform** - Works on Linux, macOS, and Windows  
✅ **NAT-Safe** - Works behind firewalls without port forwarding  
✅ **No Accounts** - No registration, no authentication servers  
✅ **Single Binary** - No dependencies, just download and run  
✅ **Terminal UI** - Interactive file browser that works everywhere

## Quick Example

### Share a folder (Terminal 1)

```bash
orb share ./documents
```

Output:

```
╔════════════════════════════════════════╗
║     Orb - Secure Folder Sharing       ║
╚════════════════════════════════════════╝

  Session:  7F9Q2A
  Passcode: 493-771
```

### Connect and browse (Terminal 2)

```bash
orb connect 7F9Q2A --passcode 493-771
```

Opens an interactive file browser where you can:

- Navigate directories
- Download files
- Browse securely over encrypted tunnel

## How It Works

```
┌─────────┐                  ┌───────────┐                  ┌─────────┐
│ Sharer  │◄────encrypted────►│   Relay   │◄────encrypted────►│Receiver │
│  (CLI)  │                  │  Server   │                  │  (CLI)  │
└─────────┘                  └───────────┘                  └─────────┘
     │                              │                              │
     └──────────────────────────────┴──────────────────────────────┘
              All encryption happens at the edges
              Relay never sees plaintext data
```

1. **Sharer** creates a session and gets a session ID + passcode
2. **Relay Server** acts as a blind byte pipe (never sees plaintext)
3. **Receiver** connects using session ID + passcode
4. **Encrypted Tunnel** established with mutual authentication
5. **Files** are accessed through secure, sandboxed filesystem

## Security Highlights

- **Argon2id** key derivation (64MB memory-hard, GPU-resistant)
- **Noise Protocol** for authenticated key exchange with perfect forward secrecy
- **ChaCha20-Poly1305** authenticated encryption
- **Path sanitization** prevents traversal attacks
- **Symlink protection** blocks escapes outside shared directory
- **Rate limiting** prevents brute force (5 attempts → lockout)
- **Session expiration** automatic cleanup after 24 hours

## Use Cases

### 🚀 Development Teams

Share build artifacts, logs, or config files between team members without uploading to cloud storage.

### 🔒 Security Auditing

Share sensitive files for security review without leaving traces on third-party servers.

### 📦 Large File Transfer

Transfer large files between machines without size limits or cloud storage costs.

### 🏢 Remote Work

Access files from home server while traveling, without VPN or port forwarding.

### 🛠️ DevOps

Share deployment artifacts or debugging information with contractors or remote teams.

## Getting Started

Choose your path:

<div class="grid cards" markdown>

- :material-clock-fast:{ .lg .middle } **Quick Start**

  ***

  Get up and running in 5 minutes

  [:octicons-arrow-right-24: Quick Start Guide](quickstart.md)

- :material-download:{ .lg .middle } **Installation**

  ***

  Install Orb on your system

  [:octicons-arrow-right-24: Install Now](getting-started/installation.md)

- :material-book-open-variant:{ .lg .middle } **User Guide**

  ***

  Learn how to use all features

  [:octicons-arrow-right-24: Read Guide](user-guide/commands.md)

- :material-shield-lock:{ .lg .middle } **Security**

  ***

  Understand the security architecture

  [:octicons-arrow-right-24: Security Details](security.md)

</div>

## Why Orb?

### No Cloud, No Problem

Traditional file sharing requires:

- ❌ Cloud accounts
- ❌ Upload time
- ❌ Storage limits
- ❌ Subscription fees
- ❌ Trust in third parties

Orb provides:

- ✅ Direct peer-to-peer sharing
- ✅ Instant access (no uploads)
- ✅ No size limits
- ✅ Free and open source
- ✅ Zero-trust architecture

### Designed for Security

Orb is built from the ground up with security in mind:

- **No long-term secrets** - Sessions expire automatically
- **No static keys** - Ephemeral keys for perfect forward secrecy
- **No plaintext leakage** - Relay server is blind to all data
- **Constant-time operations** - Resistant to timing attacks
- **Memory-hard key derivation** - Resistant to brute force

## Community & Support

- **Documentation**: You're reading it! 📚
- **GitHub**: [github.com/Zayan-Mohamed/orb](https://github.com/Zayan-Mohamed/orb)
- **Issues**: Report bugs and request features
- **Discussions**: Ask questions and share experiences

## License

Orb is open source software licensed under the [MIT License](about/license.md).

---

<div align="center">
  <p><strong>Stay safe. Stay paranoid. 🛡️</strong></p>
</div>
