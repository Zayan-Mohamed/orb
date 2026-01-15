# First Steps

Welcome! This guide will walk you through your first experience with Orb.

## Quick Overview

Orb enables secure, peer-to-peer file sharing through an encrypted tunnel. The workflow is simple:

1. **Start a Relay Server** - Acts as a blind intermediary (optional, public relay available)
2. **Share a Directory** - Creates a session and waits for connections
3. **Connect** - Uses session credentials to access shared files

## Your First File Share

### Step 1: Start Sharing

On the computer with files you want to share:

```bash
# Share the current directory
orb share .

# Or share a specific directory
orb share /path/to/folder
```

You'll see output like:

```
Session ID: abc123
Passcode: xyz789
Relay: ws://localhost:8080

Waiting for connection...
```

**Important:** Keep this terminal open and note the Session ID and Passcode!

### Step 2: Connect from Another Device

On the computer that wants to access the files:

```bash
# Connect using the session credentials
orb connect --session abc123 --passcode xyz789
```

If the relay is on a different server:

```bash
orb connect --session abc123 --passcode xyz789 --relay ws://relay.example.com:8080
```

### Step 3: Browse Files

After connecting, you'll see the TUI file browser:

```
┌─ Remote Files ────────────────────────┐
│   documents/                        │
│   photos/                           │
│   README.md                         │
│   report.pdf                        │
└───────────────────────────────────────┘
```

**Navigation:**

- `↑/↓` - Move cursor
- `Enter` - Enter directory or download file
- `Backspace` - Go to parent directory
- `q` - Quit

### Step 4: Download Files

1. Navigate to the file you want
2. Press `Enter`
3. The file downloads to your current directory
4. You'll see a success message

## Understanding the Output

### Sharing Terminal

```
[INFO] Creating session...
[INFO] Session created: abc123
[INFO] Passcode: xyz789
[INFO] Connecting to relay...
[INFO] Connected to relay
[INFO] Waiting for connection...
[INFO] Client connected! Starting handshake...
[INFO] Handshake complete
[INFO] Ready to serve files
```

### Connecting Terminal

```
[INFO] Connecting to relay...
[INFO] Connected to relay
[INFO] Starting handshake...
[INFO] Handshake complete
[INFO] Launching file browser...
```

## Common Scenarios

### Share a Specific Directory

```bash
orb share ~/Documents
```

### Use a Custom Relay

```bash
# Start your relay server
orb relay --port 9090

# Share using custom relay
orb share . --relay ws://localhost:9090

# Connect using custom relay
orb connect --session abc123 --passcode xyz789 --relay ws://localhost:9090
```

### Quick Test on One Machine

Open three terminals:

**Terminal 1 - Relay:**

```bash
orb relay
```

**Terminal 2 - Share:**

```bash
orb share ~/test-files
# Note the session ID and passcode
```

**Terminal 3 - Connect:**

```bash
orb connect --session <ID> --passcode <CODE>
```

## Security Notes

 **Your files are secure:**

- End-to-end encrypted with ChaCha20-Poly1305
- Relay server cannot see file contents
- Passcode protects session access
- Sessions expire after 24 hours

**Important:**

- Share the passcode securely (encrypted messaging, in person)
- Don't reuse session credentials
- Only share directories you trust

## What's Happening Behind the Scenes?

1. **Session Creation**: Server generates unique credentials
2. **Relay Connection**: Both peers connect to relay via WebSocket
3. **Noise Handshake**: Peers establish encrypted channel using passcode
4. **Encrypted Tunnel**: All file operations encrypted end-to-end
5. **Sandboxed Access**: File operations restricted to shared directory

## Next Steps

- Learn more about [Sharing Files](../user-guide/sharing.md)
- Explore [TUI Browser](../user-guide/tui.md) features
- Understand [Security](../security/cryptography.md) details
- Check [Troubleshooting](../user-guide/troubleshooting.md) for common issues

## Need Help?

- Run `orb --help` for command help
- Read the [FAQ](../about/faq.md)
- Check [Troubleshooting Guide](../user-guide/troubleshooting.md)
