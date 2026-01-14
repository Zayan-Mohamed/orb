# Sharing Files

Learn how to securely share files and directories using Orb.

## Overview

Sharing in Orb follows a simple workflow:

1. **Select directory** - Choose what to share
2. **Create session** - Generate credentials
3. **Share credentials** - Send to recipient securely
4. **Serve files** - Wait for connection and serve files

## Basic Sharing

### Share Current Directory

```bash
cd /path/to/files
orb share .
```

### Share Specific Directory

```bash
orb share /home/user/documents
```

### Share with Custom Relay

```bash
orb share ~/photos --relay ws://my-relay.com:8080
```

## Understanding Session Credentials

When you start sharing, Orb generates:

```
Session ID: a1b2c3d4e5f6
Passcode: secure-random-passcode-here
Relay: ws://localhost:8080
```

### Session ID

- Unique identifier for this sharing session
- Used to route connections through the relay
- Safe to share publicly (like a username)
- 12 characters, alphanumeric

### Passcode

- Secret authentication credential
- Derived to encryption keys
- **Must be kept secret**
- Random, high-entropy string
- Required for connection

### Relay URL

- WebSocket server facilitating connection
- Can be localhost, self-hosted, or public
- Must be reachable by both parties

## Sharing Best Practices

### 1. Choose What to Share Carefully

```bash
# ✅ Good: Specific directory
orb share ~/project/public-docs

# ⚠️ Risky: Entire home directory
orb share ~

# ❌ Bad: System directories
orb share /
```

### 2. Transmit Credentials Securely

**Good methods:**

- Encrypted messaging (Signal, WhatsApp)
- Password managers (shared vault)
- In person
- Phone call
- Encrypted email

**Avoid:**

- Plain email
- SMS
- Public chat
- Social media
- Shared notes

### 3. Monitor the Session

Watch the sharing terminal for:

```
[INFO] Client connected from relay
[INFO] Handshake complete
[INFO] Serving file: document.pdf
[INFO] Connection closed
```

### 4. Stop When Done

Press `Ctrl+C` to stop sharing:

```
^C
[INFO] Shutting down...
[INFO] Session terminated
```

## Advanced Sharing

### Time-Limited Sharing

```bash
# Auto-stop after 1 hour
timeout 1h orb share ~/files

# Auto-stop after 30 minutes
timeout 30m orb share ~/sensitive-docs
```

### Read-Only Sharing

Orb shares are inherently read-only. The connector can:

- ✅ List files
- ✅ Read file contents
- ✅ Download files
- ❌ Modify files
- ❌ Delete files
- ❌ Upload files

### Multiple Connections

Currently, one session = one connection. For multiple users:

```bash
# Start multiple share sessions
orb share ~/docs --relay ws://localhost:8080  # Terminal 1
orb share ~/docs --relay ws://localhost:8080  # Terminal 2

# Each gets unique credentials
```

## Directory Structure

### What Gets Shared

When you share a directory, the connector sees:

- All files in the directory
- All subdirectories
- Hidden files (dotfiles)
- Symlinks (as regular files/dirs)

### Path Isolation

Orb enforces strict path sandboxing:

```
Shared: /home/user/project
Accessible:
  ✅ /home/user/project/file.txt
  ✅ /home/user/project/subdir/file.txt
  ❌ /home/user/other-project/file.txt
  ❌ /home/user/project/../sensitive.txt
```

Protection against:

- Path traversal (`../../../etc/passwd`)
- Symlink escapes (symlinks pointing outside)
- Absolute path access

## Session Lifecycle

### 1. Creation

```
[INFO] Creating session...
[INFO] Session ID: abc123
[INFO] Passcode: xyz789
```

- Server generates random credentials
- Session stored with 24-hour expiration
- Rate limit: 5 failed attempts

### 2. Active

```
[INFO] Waiting for connection...
[INFO] Client connected!
[INFO] Handshake in progress...
[INFO] Handshake complete
[INFO] Ready to serve files
```

- Sharer waits for connector
- Performs Noise Protocol handshake
- Establishes encrypted tunnel

### 3. Serving

```
[INFO] Request: LIST /
[INFO] Request: STAT /documents
[INFO] Request: READ /documents/report.pdf
```

- Connector browses and downloads files
- All operations encrypted end-to-end
- Logs show activity

### 4. Termination

```
[INFO] Connection closed
[INFO] Session ended
```

- User presses `Ctrl+C`
- Connector disconnects
- Session can be reused if within 24h

## Troubleshooting Sharing

### Cannot Connect to Relay

```
Error: dial tcp: connection refused
```

**Solutions:**

- Check relay is running: `curl http://localhost:8080`
- Verify relay URL is correct
- Check firewall settings
- Try different port

### Session Creation Failed

```
Error: failed to create session
```

**Solutions:**

- Check session server is reachable
- Verify network connectivity
- Check server logs for errors

### Handshake Timeout

```
Error: handshake timeout
```

**Solutions:**

- Verify passcode is correct
- Check both parties using same relay
- Ensure session hasn't expired
- Try creating new session

### Permission Denied

```
Error: permission denied: /path/to/file
```

**Solutions:**

- Check file permissions
- Verify you own the files
- Run as appropriate user
- Check directory is readable

## Security Considerations

### Encryption

All file data is encrypted using:

- **ChaCha20-Poly1305** for AEAD encryption
- **Noise Protocol** for key exchange
- **Argon2id** for password-based key derivation

The relay server sees only encrypted bytes.

### Access Control

- Passcode required for connection
- Rate limiting prevents brute force
- Sessions expire after 24 hours
- No persistent authentication

### Privacy

- File names are encrypted in transit
- Directory structure is private
- Relay cannot see file metadata
- No logging of decrypted content

## Use Cases

### Temporary File Drop

```bash
# Create temporary directory
mkdir /tmp/share
cp files-to-share/* /tmp/share/

# Share
orb share /tmp/share

# After recipient downloads, cleanup
rm -rf /tmp/share
```

### Project Collaboration

```bash
# Share project directory
cd ~/projects/webapp
orb share .

# Colleague connects and downloads needed files
```

### Remote Access

```bash
# Access home files from work
ssh home-server
orb share ~/documents

# Connect from work machine
orb connect --session <ID> --passcode <CODE> --relay ws://home-server:8080
```

### Client Deliverables

```bash
# Prepare deliverables
mkdir client-delivery
cp final-*.pdf client-delivery/
cp -r assets/ client-delivery/

# Share with client
orb share client-delivery
```

## Next Steps

- Learn how to [Connect](connecting.md)
- Explore [TUI Browser](tui.md)
- Read [Security Details](../security/cryptography.md)
- Check [Troubleshooting](troubleshooting.md)
