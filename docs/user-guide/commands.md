# Command Reference

Complete reference for all Orb commands and options.

## Global Flags

These flags are available for all commands:

- `--help`, `-h` - Display help information
- `--version` - Display version information

## orb share

Share a directory securely through the relay server.

### Synopsis

```bash
orb share [directory] [flags]
```

### Arguments

- `directory` - Path to directory to share (default: current directory)

### Flags

- `--relay string` - Relay server WebSocket URL (default: "ws://localhost:8080")
- `--session-server string` - Session creation server URL (default: "http://localhost:8080")

### Description

The `share` command creates a secure sharing session for a directory:

1. Creates a session with unique ID and passcode
2. Connects to the relay server
3. Waits for incoming connections
4. Serves files from the specified directory over an encrypted tunnel

### Examples

Share current directory:

```bash
orb share .
```

Share specific directory:

```bash
orb share /home/user/documents
```

Use custom relay:

```bash
orb share ~/photos --relay ws://relay.example.com:8080
```

### Output

```
Session ID: abc123def456
Passcode: secure-random-passcode
Relay: ws://localhost:8080

Share these credentials securely with the recipient.
Waiting for connection...
```

### Security Notes

- Session credentials are printed to stdout
- Keep the terminal open while sharing
- Press `Ctrl+C` to stop sharing
- Sessions expire after 24 hours
- Only files within the shared directory are accessible

---

## orb connect

Connect to a shared directory and browse files.

### Synopsis

```bash
orb connect [flags]
```

### Flags

- `--session string` - Session ID (required)
- `--passcode string` - Session passcode (required)
- `--relay string` - Relay server WebSocket URL (default: "ws://localhost:8080")

### Description

The `connect` command establishes a connection to a shared directory:

1. Connects to the relay server
2. Performs encrypted handshake with sharer
3. Launches interactive TUI file browser
4. Allows browsing and downloading files

### Examples

Basic connection:

```bash
orb connect --session abc123 --passcode xyz789
```

Custom relay:

```bash
orb connect --session abc123 --passcode xyz789 --relay ws://relay.example.com:8080
```

### TUI Controls

Once connected, use these keys in the file browser:

- `↑` / `k` - Move cursor up
- `↓` / `j` - Move cursor down
- `Enter` - Enter directory or download file
- `Backspace` - Go to parent directory
- `q` / `Ctrl+C` - Quit browser

### Output

```
Connecting to relay ws://localhost:8080...
Connected to relay
Starting handshake...
Handshake complete
Launching file browser...

[TUI Browser Interface]
```

### Error Messages

- `Invalid session credentials` - Wrong session ID or passcode
- `Connection refused` - Relay server not reachable
- `Session expired` - Session older than 24 hours
- `Handshake failed` - Encryption handshake failed

---

## orb relay

Start a relay server to facilitate connections.

### Synopsis

```bash
orb relay [flags]
```

### Flags

- `--host string` - Host to bind to (default: "localhost")
- `--port int` - Port to listen on (default: 8080)

### Description

The `relay` command starts a WebSocket server that acts as a blind intermediary:

1. Accepts WebSocket connections from both sharers and connectors
2. Pairs connections based on session ID
3. Forwards encrypted messages between peers
4. Cannot decrypt or inspect traffic

### Examples

Start relay on default port:

```bash
orb relay
```

Start on custom port:

```bash
orb relay --port 9090
```

Bind to all interfaces (for public relay):

```bash
orb relay --host 0.0.0.0 --port 8080
```

### Output

```
Starting Orb relay server...
Listening on localhost:8080
WebSocket endpoint: ws://localhost:8080/ws

[INFO] Accepting connections...
```

### Session Management

The relay server provides HTTP endpoints:

- `POST /sessions` - Create new session
  - Returns: Session ID and passcode
  - Body: None required
- `GET /sessions/{id}` - Verify session exists
  - Returns: Session status
  - Requires: Session ID in URL

### WebSocket Protocol

- `GET /ws?session={id}` - Establish WebSocket connection
  - Requires: Session ID in query parameter
  - Pairs connections with matching session ID
  - Forwards all frames between paired connections

### Deployment Notes

- Production relay should use TLS (`wss://`)
- Implement rate limiting for production use
- Monitor connection counts and bandwidth
- Set up logging for security audits

---

## orb help

Display help information for any command.

### Synopsis

```bash
orb help [command]
```

### Examples

General help:

```bash
orb help
```

Help for specific command:

```bash
orb help share
orb help connect
orb help relay
```

---

## orb version

Display version information.

### Synopsis

```bash
orb version
```

### Output

```
Orb v1.0.0
Build: 2024-01-15
Go version: go1.21.0
```

---

## Environment Variables

Orb respects the following environment variables:

### ORB_RELAY

Default relay server URL:

```bash
export ORB_RELAY="ws://relay.example.com:8080"
orb share ~/files  # Uses ORB_RELAY
```

### ORB_SESSION_SERVER

Default session server URL:

```bash
export ORB_SESSION_SERVER="http://relay.example.com:8080"
orb share ~/files  # Uses ORB_SESSION_SERVER
```

### ORB_DEBUG

Enable debug logging:

```bash
export ORB_DEBUG=1
orb share ~/files  # Prints debug messages
```

---

## Exit Codes

Orb uses the following exit codes:

- `0` - Success
- `1` - General error
- `2` - Connection error
- `3` - Authentication error
- `4` - File system error
- `130` - Interrupted by user (Ctrl+C)

---

## Configuration Files

Orb currently does not use configuration files. All options must be specified via command-line flags or environment variables.

Future versions may support:

- `~/.config/orb/config.yaml` - User configuration
- `~/.orb/known_relays` - List of trusted relays

---

## Shell Completion

Generate shell completion scripts:

### Bash

```bash
orb completion bash > /etc/bash_completion.d/orb
```

### Zsh

```bash
orb completion zsh > /usr/local/share/zsh/site-functions/_orb
```

### Fish

```bash
orb completion fish > ~/.config/fish/completions/orb.fish
```

### PowerShell

```powershell
orb completion powershell > orb.ps1
```

---

## Advanced Usage

### Chaining Commands

```bash
# Share and log session info
orb share ~/files 2>&1 | tee session.log

# Connect with timeout
timeout 5m orb connect --session abc123 --passcode xyz789
```

### Background Sharing

```bash
# Share in background (use with caution)
orb share ~/files > session.log 2>&1 &

# Check status
jobs

# Bring to foreground
fg
```

### Multiple Sessions

```bash
# Share multiple directories simultaneously
orb share ~/docs --relay ws://localhost:8080 &
orb share ~/photos --relay ws://localhost:8080 &

# Each gets unique session ID
```

---

## Next Steps

- Learn about [Sharing](sharing.md) in detail
- Explore [TUI Browser](tui.md) features
- Check [Troubleshooting](troubleshooting.md) guide
