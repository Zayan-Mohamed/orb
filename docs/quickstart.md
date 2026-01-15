# Orb Quick Start Guide

## Installation

### Option 1: Build from Source

```bash
# Clone and build
cd /home/zayan/Documents/myProjects/orb
go build -o orb .

# Or use make
make build-local

# Move to PATH
sudo mv orb /usr/local/bin/
```

### Option 2: Cross-Platform Build

```bash
# Build for all platforms
./build.sh

# Binaries will be in build/ directory
ls build/
```

## Quick Demo (localhost)

### Terminal 1: Start Relay

```bash
orb relay --listen :8080
```

Expected output:

```
Starting Orb relay server...
Listening on :8080

Security notes:
  • The relay server never sees plaintext data
  • All encryption happens at the edges
  • Sessions expire automatically
```

### Terminal 2: Share a Folder

```bash
# Create test folder
mkdir -p /tmp/orb-test
echo "Hello from Orb!" > /tmp/orb-test/test.txt

# Share it
orb share /tmp/orb-test --relay http://localhost:8080
```

Expected output:

```
╔════════════════════════════════════════╗
║     Orb - Secure Folder Sharing       ║
╚════════════════════════════════════════╝

  Session:  ABC123
  Passcode: 123-456

Share these credentials with the receiver.
Waiting for connection...
```

### Terminal 3: Connect to Share

```bash
orb connect ABC123 --passcode 123-456 --relay http://localhost:8080 --tui
```

You'll see an interactive file browser. Use:

- **↑/↓**: Navigate files
- **Enter**: Open directory or download file
- **Backspace**: Go to parent directory
- **d**: Download selected file
- **q**: Quit

## Real-World Usage

### Share with Remote User

1. **Setup** (once):

   ```bash
   # Deploy relay on VPS
   ssh user@your-vps.com
   ./orb relay --listen 0.0.0.0:8080
   ```

2. **Share** (from your machine):

   ```bash
   orb share ./my-documents \
     --relay http://your-vps.com:8080 \
     --readonly
   ```

3. **Send Credentials** (securely):

   - Send session ID and passcode via Signal/WhatsApp
   - **Never** via email or SMS!

4. **Receiver Connects**:
   ```bash
   orb connect SESSION_ID \
     --passcode PASSCODE \
     --relay http://your-vps.com:8080
   ```

## Common Scenarios

### Scenario 1: Share Read-Only

```bash
orb share ./project-docs --readonly
```

Files can be viewed and downloaded, but not modified or deleted.

### Scenario 2: Collaborative Editing

```bash
# Sharer (read-write)
orb share ./shared-workspace

# Receiver can upload/modify files
# Use with trusted parties only!
```

### Scenario 3: Quick File Transfer

```bash
# Share
orb share ./large-video.mp4

# Receiver downloads via TUI
# Session auto-expires after transfer
```

## Keyboard Shortcuts (TUI Mode)

| Key         | Action                       |
| ----------- | ---------------------------- |
| ↑/↓         | Navigate up/down             |
| Enter       | Open folder or download file |
| Backspace   | Parent directory             |
| d           | Download selected file       |
| /           | Search/filter files          |
| q or Ctrl+C | Quit                         |

## Troubleshooting

### "Failed to connect to relay"

**Cause**: Relay server not running or wrong URL

**Fix**:

```bash
# Check relay is running
curl http://localhost:8080/health

# Or start relay
orb relay --listen :8080
```

### "Authentication failed"

**Cause**: Wrong passcode or locked session

**Fix**:

- Verify passcode is correct
- Check for typos
- Session locks after 5 failed attempts
- Create new session if locked

### "Session expired"

**Cause**: Session older than 24 hours

**Fix**:

- Create new session
- Sessions auto-expire for security

### "Failed to read file"

**Cause**: Permission issues or file deleted

**Fix**:

- Check file still exists
- Verify sharer has read permission
- Check if share is read-only mode

## Security Tips

1. **Strong Passcodes**

   - Let Orb generate random passcodes
   - Don't reuse passcodes
   - Share via secure channels only

2. **Read-Only by Default**

   - Use `--readonly` unless write access needed
   - Limits damage if receiver is compromised

3. **Session Hygiene**

   - Create new session for each transfer
   - Verify receiver identity before sharing credentials
   - Sessions auto-expire after 24 hours

4. **Network Security**

   - Use HTTPS/WSS for relay in production
   - Deploy relay behind reverse proxy with TLS
   - Consider VPN for extra security

5. **Data Sensitivity**
   - Orb is secure, but verify relay server trust
   - For highly sensitive data, run your own relay
   - Don't share credentials in plaintext

## Performance Tips

1. **Large Files**

   - Files download in chunks
   - Close other applications if memory constrained
   - Use stable network connection

2. **Multiple Files**

   - Download files individually via TUI
   - Or batch download by navigating directories

3. **Slow Networks**
   - Orb handles poor connections gracefully
   - Automatic reconnection on disconnect
   - Keep-alive prevents timeout on idle

## Advanced Usage

### Custom Relay Port

```bash
orb relay --listen :9090
orb share ./folder --relay http://localhost:9090
```

### Behind Firewall

Both peers can be behind NAT/firewall - relay handles traversal:

```bash
# Both peers behind NAT - still works!
# Relay must be on public internet
```

### Multiple Sessions

You can run multiple share sessions simultaneously:

```bash
# Terminal 1
orb share ./docs --relay http://relay:8080

# Terminal 2
orb share ./photos --relay http://relay:8080

# Different session IDs, isolated sessions
```

## Monitoring

### Relay Server Logs

Relay logs connections (no content):

```
Session created: ABC123
Sharer connected: session=ABC123
Receiver connected: session=ABC123
Session closed: ABC123
```

**Note**: File names and content are never logged!

### Connection Status

Check if tunnel is alive:

```bash
# From receiver side, connection status shown in TUI
# If tunnel dies, TUI shows error
```

## Getting Help

### Command Help

```bash
orb --help
orb share --help
orb connect --help
orb relay --help
```

### Version Info

```bash
orb version
```

### Debug Mode

```bash
# More verbose output
orb share ./folder --relay http://localhost:8080 --debug
```

## What's Next?

- Read the [full documentation](index.md)
- Review [Security Overview](security.md) for security details
- Check [Architecture](architecture.md) for system design
- Report issues on GitHub

---

**Note**: Orb is designed for security. Follow best practices for secure file sharing.
