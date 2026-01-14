# Usage Examples

This page provides practical examples for common Orb use cases.

## Basic Examples

### 1. Share a Directory

Share the current directory:

```bash
orb share .
```

Share a specific directory:

```bash
orb share /home/user/documents
```

Share with a custom relay:

```bash
orb share ~/photos --relay ws://relay.example.com:8080
```

### 2. Connect to a Share

Basic connection:

```bash
orb connect --session abc123 --passcode xyz789
```

Connection with custom relay:

```bash
orb connect --session abc123 --passcode xyz789 --relay ws://my-relay.com:8080
```

### 3. Run a Relay Server

Default configuration:

```bash
orb relay
```

Custom port:

```bash
orb relay --port 9090
```

Bind to specific interface:

```bash
orb relay --host 0.0.0.0 --port 8080
```

## Advanced Examples

### Self-Hosted Relay

**Server Setup:**

```bash
# Start relay on server
ssh user@relay.example.com
orb relay --host 0.0.0.0 --port 8080
```

**Share files:**

```bash
orb share ~/project --relay ws://relay.example.com:8080
```

**Connect:**

```bash
orb connect --session <ID> --passcode <CODE> --relay ws://relay.example.com:8080
```

### Temporary File Drop

Quick file sharing that expires after one connection:

```bash
# Share a file temporarily
orb share /tmp/transfer

# After download, press Ctrl+C to stop sharing
```

### Project Collaboration

Share a project directory with a team member:

```bash
# Developer A shares project
cd ~/projects/myapp
orb share .

# Developer B connects and downloads specific files
orb connect --session <ID> --passcode <CODE>
# Navigate to needed files in TUI
# Press Enter to download
```

### Large File Transfer

Transfer large files between machines:

```bash
# Machine A: Share directory with large files
orb share ~/large-files

# Machine B: Connect and download
orb connect --session <ID> --passcode <CODE>
# Files are downloaded to current directory
```

## Use Case Scenarios

### 1. Remote Work

**Scenario:** Access work files from home

```bash
# At office
orb share ~/work-documents

# At home (note the credentials)
orb connect --session <ID> --passcode <CODE>
```

### 2. Client Deliverables

**Scenario:** Send files to a client securely

```bash
# Create deliverables directory
mkdir client-deliverables
cp final-report.pdf client-deliverables/
cp screenshots/* client-deliverables/

# Share it
orb share client-deliverables

# Send credentials to client via secure channel
# Client connects and downloads
```

### 3. Emergency File Access

**Scenario:** Need a file urgently from another machine

```bash
# On machine with files (via SSH)
ssh home-server
cd /important/files
orb share .

# On current machine
orb connect --session <ID> --passcode <CODE>
# Download the needed file
```

### 4. Cross-Platform Transfer

**Scenario:** Move files between Windows, Mac, and Linux

```bash
# Windows
orb.exe share C:\Users\Alice\Documents

# macOS
orb connect --session <ID> --passcode <CODE>
# Downloads work seamlessly across platforms
```

### 5. Development Environment Setup

**Scenario:** Share config files with team

```bash
# Share dotfiles and configs
orb share ~/.config

# Team member downloads specific configs
# .vimrc, .bashrc, etc.
```

## Scripting Examples

### Automated Sharing Script

```bash
#!/bin/bash
# share-script.sh

DIRECTORY=$1
RELAY="ws://relay.example.com:8080"

echo "Starting share for: $DIRECTORY"
orb share "$DIRECTORY" --relay "$RELAY"
```

Usage:

```bash
chmod +x share-script.sh
./share-script.sh ~/documents
```

### Batch Connection Script

```bash
#!/bin/bash
# connect-script.sh

SESSION=$1
PASSCODE=$2
RELAY=${3:-ws://localhost:8080}

echo "Connecting to session: $SESSION"
orb connect --session "$SESSION" --passcode "$PASSCODE" --relay "$RELAY"
```

Usage:

```bash
chmod +x connect-script.sh
./connect-script.sh abc123 xyz789
```

### Session Logger

```bash
#!/bin/bash
# share-and-log.sh

LOGFILE="orb-sessions.log"
DIRECTORY=$1

echo "=== New Share Session ===" >> "$LOGFILE"
echo "Date: $(date)" >> "$LOGFILE"
echo "Directory: $DIRECTORY" >> "$LOGFILE"

orb share "$DIRECTORY" 2>&1 | tee -a "$LOGFILE"
```

## Integration Examples

### With SSH

```bash
# Share via SSH tunnel
ssh -L 8080:localhost:8080 user@remote-server
# In another terminal
orb share ~/files --relay ws://localhost:8080
```

### With systemd

Create `/etc/systemd/system/orb-relay.service`:

```ini
[Unit]
Description=Orb Relay Server
After=network.target

[Service]
Type=simple
User=orb
ExecStart=/usr/local/bin/orb relay --host 0.0.0.0 --port 8080
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable orb-relay
sudo systemctl start orb-relay
```

### With Docker

See [Docker Deployment](../deployment/docker.md) for containerized examples.

## Tips and Tricks

### 1. Alias for Quick Sharing

Add to `.bashrc` or `.zshrc`:

```bash
alias orbshare='orb share .'
alias orbconnect='orb connect'
```

### 2. Environment Variables

Set default relay:

```bash
export ORB_RELAY="ws://relay.example.com:8080"
orb share ~/files
```

### 3. Quick Directory Share

```bash
# Share and go back to work
orb share ~/transfer &
# When done, bring to foreground and stop
fg
# Ctrl+C
```

### 4. Check Relay Status

```bash
# Test relay connectivity
curl -I http://relay.example.com:8080
```

## Common Patterns

### Pattern 1: Secure Handoff

```bash
# Step 1: Create temporary directory
mkdir /tmp/handoff
cp secret-file.txt /tmp/handoff/

# Step 2: Share temporarily
orb share /tmp/handoff

# Step 3: Send credentials securely
# Step 4: After download, cleanup
rm -rf /tmp/handoff
```

### Pattern 2: Rolling Deployment

```bash
# Deploy files to multiple servers
for server in server1 server2 server3; do
    ssh $server "orb connect --session <ID> --passcode <CODE>"
done
```

### Pattern 3: Backup Retrieval

```bash
# Share backup directory
orb share /backups/latest

# Remote retrieval
orb connect --session <ID> --passcode <CODE>
# Download needed files
```

## Next Steps

- Learn about [Command Reference](../user-guide/commands.md)
- Explore [TUI Features](../user-guide/tui.md)
- Check [Troubleshooting](../user-guide/troubleshooting.md)
