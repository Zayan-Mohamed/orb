# Troubleshooting

Common issues and solutions for Orb.

## Connection Issues

### Cannot Connect to Relay

**Error:**

```
Error: dial tcp [::1]:8080: connect: connection refused
```

**Causes:**

- Relay server not running
- Wrong relay address
- Firewall blocking connection
- Network unavailable

**Solutions:**

1. **Verify relay is running:**

   ```bash
   curl -I http://localhost:8080
   # Should return HTTP response
   ```

2. **Check relay address:**

   ```bash
   # Ensure protocol and port are correct
   ws://localhost:8080  # Correct
   wss://localhost:8080 # TLS version
   http://localhost:8080 # Wrong (HTTP, not WebSocket)
   ```

3. **Test with local relay:**

   ```bash
   # Terminal 1: Start relay
   orb relay

   # Terminal 2: Try connection
   orb connect --session <ID> --passcode <CODE> --relay ws://localhost:8080
   ```

4. **Check firewall:**

   ```bash
   # Linux
   sudo ufw allow 8080

   # macOS
   # System Preferences → Security → Firewall

   # Windows
   # Windows Defender Firewall → Allow an app
   ```

### Session Not Found

**Error:**

```
Error: session not found
```

**Causes:**

- Wrong session ID
- Session expired (>24h)
- Session never created
- Typo in credentials

**Solutions:**

1. **Verify session ID:**

   ```bash
   # Double-check session ID from sharer
   # It's case-sensitive: abc123 ≠ ABC123
   ```

2. **Check session age:**

   ```bash
   # Sessions expire after 24 hours
   # Ask sharer to create new session
   ```

3. **Create fresh session:**
   ```bash
   # Sharer creates new session
   orb share ~/files
   # Note new credentials
   ```

### Authentication Failed

**Error:**

```
Error: authentication failed
Error: handshake failed
```

**Causes:**

- Incorrect passcode
- Passcode has extra spaces
- Wrong key derivation
- Crypto mismatch

**Solutions:**

1. **Verify passcode exactly:**

   ```bash
   # Passcode is case-sensitive
   # Check for spaces: "pass code" vs "passcode"
   # Copy-paste to avoid typos
   ```

2. **Request credentials again:**

   ```bash
   # Ask sharer to resend passcode
   # Use secure channel
   ```

3. **Create new session:**
   ```bash
   # If passcode lost/wrong, start over
   orb share ~/files  # New credentials
   ```

### Handshake Timeout

**Error:**

```
Error: context deadline exceeded
Error: handshake timeout
```

**Causes:**

- Network latency
- Sharer not connected
- Firewall blocking packets
- Relay issues

**Solutions:**

1. **Verify sharer is connected:**

   ```bash
   # Check sharer terminal shows:
   # "Waiting for connection..."
   ```

2. **Test network speed:**

   ```bash
   ping -c 5 relay.example.com
   # Check latency
   ```

3. **Use closer relay:**

   ```bash
   # Self-host relay geographically closer
   # Or use relay with better connectivity
   ```

4. **Try again:**
   ```bash
   # Temporary network issue
   # Simply retry connection
   ```

## Sharing Issues

### Permission Denied

**Error:**

```
Error: permission denied
```

**Causes:**

- No read access to directory
- Files owned by different user
- SELinux/AppArmor blocking

**Solutions:**

1. **Check permissions:**

   ```bash
   ls -la /path/to/share
   # Ensure files are readable
   ```

2. **Fix permissions:**

   ```bash
   chmod -R +r /path/to/share
   # Make all files readable
   ```

3. **Run as correct user:**
   ```bash
   # If files owned by different user
   sudo -u owner orb share /path/to/files
   ```

### Directory Not Found

**Error:**

```
Error: directory does not exist
```

**Causes:**

- Typo in path
- Directory moved/deleted
- Relative vs absolute path

**Solutions:**

1. **Verify path:**

   ```bash
   ls /path/to/directory
   # Confirm exists
   ```

2. **Use absolute path:**

   ```bash
   # Instead of: orb share ../files
   orb share /home/user/files
   ```

3. **Check current directory:**
   ```bash
   pwd
   # Verify you're where you think you are
   ```

## File Browser Issues

### Browser Won't Launch

**Error:**

```
Error: failed to initialize browser
```

**Causes:**

- Terminal not supported
- TTY not available
- Connection failed before launch

**Solutions:**

1. **Use supported terminal:**

   - ✅ iTerm2 (macOS)
   - ✅ Terminal.app (macOS)
   - ✅ Windows Terminal
   - ✅ Alacritty
   - ✅ GNOME Terminal
   - ❌ Dumb terminal
   - ❌ Non-interactive shell

2. **Check TTY:**

   ```bash
   tty
   # Should output /dev/pts/0 or similar
   # Not "not a tty"
   ```

3. **Verify connection:**
   ```bash
   # Ensure handshake completed
   # Check connection logs
   ```

### Cannot Download Files

**Error:**

```
Error: download failed
```

**Causes:**

- Disk full
- Permission denied (local)
- Network interruption
- File locked on sharer side

**Solutions:**

1. **Check disk space:**

   ```bash
   df -h .
   # Ensure enough space
   ```

2. **Check write permissions:**

   ```bash
   ls -la .
   # Ensure current directory is writable
   ```

3. **Change download location:**

   ```bash
   cd ~/Downloads
   orb connect --session <ID> --passcode <CODE>
   ```

4. **Try smaller file first:**
   ```bash
   # Test with small file
   # If works, issue is large file handling
   ```

### Display Corruption

**Symptoms:**

- Garbled text
- Broken boxes
- Missing characters
- Weird colors

**Causes:**

- Terminal encoding wrong
- Unicode not supported
- Terminal too small
- Color scheme issues

**Solutions:**

1. **Set UTF-8 encoding:**

   ```bash
   export LC_ALL=en_US.UTF-8
   export LANG=en_US.UTF-8
   ```

2. **Resize terminal:**

   ```bash
   # Minimum 80x24
   # Recommended 120x30
   ```

3. **Reset terminal:**

   ```bash
   reset
   # Clear any corruption
   ```

4. **Try different terminal:**
   ```bash
   # Use modern terminal emulator
   # Enable UTF-8 support
   ```

## Build Issues

### Build Fails

**Error:**

```
Error: build failed
```

**Causes:**

- Go version too old
- Missing dependencies
- Network issues downloading modules
- Platform not supported

**Solutions:**

1. **Check Go version:**

   ```bash
   go version
   # Need Go 1.21 or higher
   ```

2. **Update dependencies:**

   ```bash
   go mod tidy
   go mod download
   ```

3. **Clear cache:**

   ```bash
   go clean -cache
   go clean -modcache
   ```

4. **Build with verbose output:**
   ```bash
   go build -v -x
   # See detailed build steps
   ```

### Binary Not Found

**Error:**

```
bash: orb: command not found
```

**Causes:**

- Binary not in PATH
- Binary name incorrect
- Wrong directory

**Solutions:**

1. **Check binary location:**

   ```bash
   which orb
   # If empty, not in PATH
   ```

2. **Run with full path:**

   ```bash
   /usr/local/bin/orb --version
   ./orb --version
   ```

3. **Add to PATH:**
   ```bash
   export PATH=$PATH:/path/to/orb
   # Or copy to /usr/local/bin
   ```

## Performance Issues

### Slow File Listing

**Symptoms:**

- Directory loading takes forever
- Browser appears frozen

**Causes:**

- Large directory (thousands of files)
- Network latency
- Slow disk on sharer side

**Solutions:**

1. **Wait for completion:**

   ```bash
   # Large directories take time
   # Be patient
   ```

2. **Share smaller directory:**

   ```bash
   # Instead of sharing entire home:
   orb share ~/specific-project
   ```

3. **Check network:**
   ```bash
   ping -c 10 relay.example.com
   # Look for packet loss
   ```

### Slow Downloads

**Symptoms:**

- File download takes very long
- Much slower than expected

**Causes:**

- Network bandwidth limited
- Large file size
- Relay bottleneck
- Encryption overhead

**Solutions:**

1. **Check network speed:**

   ```bash
   # Run speed test
   # Compare to expected bandwidth
   ```

2. **Use local relay:**

   ```bash
   # If over internet, use local relay
   # Reduces hops
   ```

3. **Compress files:**

   ```bash
   # Sharer compresses before sharing
   tar -czf archive.tar.gz files/
   orb share .
   ```

4. **Split large files:**
   ```bash
   # Split large files
   split -b 100M largefile.dat part-
   # Download parts separately
   ```

## Relay Server Issues

### Relay Won't Start

**Error:**

```
Error: address already in use
```

**Causes:**

- Port 8080 already used
- Another orb relay running
- Other service using port

**Solutions:**

1. **Check port usage:**

   ```bash
   # Linux/macOS
   lsof -i :8080
   netstat -an | grep 8080

   # Windows
   netstat -ano | findstr :8080
   ```

2. **Kill conflicting process:**

   ```bash
   # Find PID from above
   kill <PID>
   ```

3. **Use different port:**
   ```bash
   orb relay --port 9090
   ```

### Relay Crashes

**Symptoms:**

- Relay exits unexpectedly
- Connections drop

**Causes:**

- Out of memory
- Too many connections
- Bug in code
- System resource limits

**Solutions:**

1. **Check logs:**

   ```bash
   # Look for error messages
   # Check system logs
   ```

2. **Increase limits:**

   ```bash
   # Linux
   ulimit -n 4096  # File descriptors
   ```

3. **Restart relay:**

   ```bash
   orb relay --port 8080
   ```

4. **Monitor resources:**
   ```bash
   # Watch memory/CPU
   top
   htop
   ```

## Platform-Specific Issues

### Windows

**Issue: Antivirus blocks Orb**

Solution:

```
1. Add orb.exe to exclusions
2. Windows Defender → Virus & threat protection
3. Add exclusion for orb.exe
```

**Issue: WebSocket connection fails**

Solution:

```powershell
# Check Windows Firewall
New-NetFirewallRule -DisplayName "Orb" -Direction Inbound -Port 8080 -Protocol TCP -Action Allow
```

### macOS

**Issue: "orb" cannot be opened because the developer cannot be verified**

Solution:

```bash
# Remove quarantine
xattr -d com.apple.quarantine orb

# Or allow in System Preferences
# Security & Privacy → Open Anyway
```

**Issue: Permission denied**

Solution:

```bash
# Make executable
chmod +x orb

# Move to /usr/local/bin
sudo mv orb /usr/local/bin/
```

### Linux

**Issue: SELinux blocks connections**

Solution:

```bash
# Temporarily disable
sudo setenforce 0

# Or create policy
# Check audit logs for denials
```

**Issue: systemd service won't start**

Solution:

```bash
# Check status
systemctl status orb-relay

# View logs
journalctl -u orb-relay -f

# Test manually first
/usr/local/bin/orb relay
```

## Getting Help

### Collect Debug Information

```bash
# Version
orb --version

# System info
uname -a  # Linux/macOS
systeminfo  # Windows

# Network
ip addr  # Linux
ifconfig  # macOS
ipconfig  # Windows
```

### Enable Debug Logging

```bash
export ORB_DEBUG=1
orb share ~/files
# More detailed output
```

### Report Issues

When reporting bugs, include:

- Orb version (`orb --version`)
- Operating system and version
- Command that failed
- Complete error message
- Steps to reproduce

### Community Support

- GitHub Issues: Report bugs
- Discussions: Ask questions
- Documentation: Check this guide
- Examples: See usage examples

## Common Error Reference

| Error                       | Meaning                   | Solution               |
| --------------------------- | ------------------------- | ---------------------- |
| `connection refused`        | Relay not reachable       | Check relay running    |
| `session not found`         | Invalid/expired session   | Create new session     |
| `authentication failed`     | Wrong passcode            | Verify credentials     |
| `permission denied`         | No file access            | Check permissions      |
| `address in use`            | Port already used         | Use different port     |
| `handshake timeout`         | Handshake didn't complete | Check network          |
| `context deadline exceeded` | Operation timed out       | Retry or check network |

## Next Steps

- Review [Commands](commands.md) reference
- Check [FAQ](../about/faq.md)
- Read [Security](../security/cryptography.md) details
