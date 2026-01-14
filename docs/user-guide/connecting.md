# Connecting to Shares

Learn how to connect to shared directories and access files securely.

## Overview

Connecting to a share requires:

1. **Session credentials** - ID and passcode from sharer
2. **Relay URL** - WebSocket server address
3. **Network connectivity** - Access to relay server

## Basic Connection

### Using Session Credentials

```bash
orb connect --session abc123 --passcode xyz789
```

### With Custom Relay

```bash
orb connect --session abc123 --passcode xyz789 --relay ws://relay.example.com:8080
```

## Connection Process

### 1. Relay Connection

```
[INFO] Connecting to relay ws://localhost:8080...
[INFO] Connected to relay
```

The connector:

- Opens WebSocket connection
- Sends session ID
- Waits to be paired with sharer

### 2. Handshake

```
[INFO] Starting handshake...
[INFO] Handshake complete
```

The connector:

- Initiates Noise Protocol handshake
- Derives encryption keys from passcode
- Establishes encrypted tunnel

### 3. File Browser Launch

```
[INFO] Launching file browser...

┌─ Remote Files ─────────────────┐
│  📁 documents/                 │
│  📁 photos/                    │
│  📄 README.md                  │
└────────────────────────────────┘
```

The TUI browser opens, showing remote files.

## Obtaining Credentials

### From Sharer

The sharer will provide:

```
Session ID: a1b2c3d4e5f6
Passcode: secure-random-passcode
Relay: ws://localhost:8080
```

### Important Notes

- **Session ID**: Can be shared openly
- **Passcode**: Must be kept secret
- **Relay**: Must be network-accessible
- **Expiration**: Sessions valid for 24 hours

## Connection Scenarios

### Local Network

Both on same network:

```bash
# Sharer starts relay
orb relay

# Sharer shares
orb share ~/files

# Connector connects
orb connect --session <ID> --passcode <CODE>
```

### Over Internet

Using public/self-hosted relay:

```bash
# Sharer
orb share ~/files --relay ws://relay.example.com:8080

# Connector (from anywhere)
orb connect --session <ID> --passcode <CODE> --relay ws://relay.example.com:8080
```

### Through SSH Tunnel

Secure connection via SSH:

```bash
# Create SSH tunnel
ssh -L 8080:localhost:8080 user@remote-server

# Connect through tunnel
orb connect --session <ID> --passcode <CODE> --relay ws://localhost:8080
```

### Behind Firewall

If connector is behind restrictive firewall:

```bash
# Ensure relay uses standard ports (80/443)
# Or configure firewall to allow WebSocket
# Use wss:// (encrypted WebSocket) for HTTPS networks
```

## Using the File Browser

Once connected, navigate with:

### Navigation Keys

- `↑` or `k` - Move up
- `↓` or `j` - Move down
- `Enter` - Enter directory or download file
- `Backspace` - Go to parent directory
- `q` - Quit

### File Operations

**Browse Directory:**

```
1. Navigate to directory
2. Press Enter
3. See contents
```

**Download File:**

```
1. Navigate to file
2. Press Enter
3. File downloads to current directory
```

**Parent Directory:**

```
1. Press Backspace
2. Go up one level
```

### Download Location

Files download to your current working directory:

```bash
cd ~/Downloads
orb connect --session <ID> --passcode <CODE>
# Files download to ~/Downloads
```

## Connection Management

### Active Connection

While connected:

```
- File browser is active
- Can navigate and download
- Connection stays open
- Encrypted throughout
```

### Disconnection

To disconnect:

```
1. Press 'q' in browser
2. Or press Ctrl+C
3. Connection terminates
```

### Reconnection

Within 24 hours:

```bash
# Use same credentials
orb connect --session <same-ID> --passcode <same-CODE>

# Fresh connection established
# File browser reopens
```

## Troubleshooting Connection

### Invalid Credentials

```
Error: authentication failed
```

**Causes:**

- Wrong session ID
- Incorrect passcode
- Typo in credentials

**Solutions:**

- Double-check session ID
- Verify passcode (case-sensitive)
- Request credentials again

### Connection Refused

```
Error: connection refused
```

**Causes:**

- Relay not running
- Wrong relay URL
- Network unreachable
- Firewall blocking

**Solutions:**

- Verify relay address
- Check relay is running: `curl http://relay:8080`
- Test network connectivity
- Check firewall rules

### Session Not Found

```
Error: session not found
```

**Causes:**

- Session expired (>24h)
- Wrong session ID
- Session not created yet

**Solutions:**

- Verify session ID
- Check session age
- Ask sharer to create new session

### Handshake Failed

```
Error: handshake timeout
```

**Causes:**

- Incorrect passcode
- Network interruption
- Crypto mismatch

**Solutions:**

- Verify passcode is exact
- Check network stability
- Try new session

### Timeout

```
Error: context deadline exceeded
```

**Causes:**

- Sharer not connected
- Network latency
- Relay issues

**Solutions:**

- Confirm sharer is online
- Check network speed
- Try different relay

## Security Considerations

### Passcode Protection

**The passcode is crucial:**

- Used to derive encryption keys
- Provides authentication
- Must remain secret

**Never share passcode via:**

- Plain email
- SMS
- Public chat
- Unencrypted channels

**Safe methods:**

- Encrypted messaging (Signal, WhatsApp)
- Phone call
- In person
- Password manager

### Verification

**Before downloading sensitive files:**

1. Verify sharer identity out-of-band
2. Confirm session credentials directly
3. Check file integrity if needed

### Network Security

**When possible:**

- Use `wss://` (encrypted) instead of `ws://`
- Connect over VPN
- Use trusted relays
- Verify relay SSL certificate

## Advanced Connection

### Environment Variables

Set default relay:

```bash
export ORB_RELAY="ws://relay.example.com:8080"
orb connect --session <ID> --passcode <CODE>
# Uses ORB_RELAY automatically
```

### Scripted Connection

```bash
#!/bin/bash
# connect.sh

SESSION=$1
PASSCODE=$2

if [ -z "$SESSION" ] || [ -z "$PASSCODE" ]; then
    echo "Usage: ./connect.sh <session> <passcode>"
    exit 1
fi

orb connect --session "$SESSION" --passcode "$PASSCODE" --relay ws://relay.example.com:8080
```

### Timeout Handling

```bash
# Timeout after 5 minutes
timeout 5m orb connect --session <ID> --passcode <CODE>
```

## Connection Lifecycle

### Phase 1: Relay Connection

```
Connecting to relay → WebSocket open → Send session ID → Paired
```

### Phase 2: Handshake

```
Derive keys → Exchange handshake → Verify → Establish tunnel
```

### Phase 3: Active Session

```
Browse files → Send requests → Receive responses → Download files
```

### Phase 4: Termination

```
User quit → Close connection → Cleanup
```

## Best Practices

### 1. Verify Before Connecting

```bash
# Check relay is reachable
curl -I http://relay.example.com:8080

# Then connect
orb connect --session <ID> --passcode <CODE> --relay ws://relay.example.com:8080
```

### 2. Use Clean Download Directory

```bash
# Create temporary download location
mkdir -p ~/Downloads/orb-session
cd ~/Downloads/orb-session
orb connect --session <ID> --passcode <CODE>
```

### 3. Close When Done

```bash
# Don't leave connections open unnecessarily
# Press 'q' to quit browser
```

### 4. One Connection at a Time

```bash
# Currently: one session = one connection
# For multiple connectors: sharer creates multiple sessions
```

## Use Cases

### Download Project Files

```bash
cd ~/projects
orb connect --session <ID> --passcode <CODE>
# Navigate to project files
# Download needed files
```

### Retrieve Documents

```bash
cd ~/Documents
orb connect --session <ID> --passcode <CODE>
# Browse document hierarchy
# Download specific documents
```

### Access Remote Files

```bash
# Connect to home server from work
orb connect --session <ID> --passcode <CODE> --relay ws://home.example.com:8080
# Download needed files
```

### Receive Client Files

```bash
# Client shares deliverables
cd ~/client-files
orb connect --session <ID> --passcode <CODE>
# Download all deliverables
```

## Next Steps

- Learn [TUI Browser](tui.md) features
- Read [Security](../security/cryptography.md) details
- Check [Troubleshooting](troubleshooting.md) guide
- Explore [Command Reference](commands.md)
