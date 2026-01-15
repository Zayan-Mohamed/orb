# TUI File Browser

Learn how to use the interactive terminal file browser.

## Overview

The TUI (Terminal User Interface) browser provides an interactive way to:

- Browse remote directory structure
- View file listings
- Download files to local machine
- Navigate directories

## Launching the Browser

The browser launches automatically after successful connection:

```bash
orb connect --session abc123 --passcode xyz789
```

## Interface Layout

```
┌─ Remote Files ─────────────────────────────────────┐
│ Current: /documents/projects                       │
├────────────────────────────────────────────────────┤
│   subfolder/                                     │
│   report.pdf                                     │
│   data.xlsx                                      │
│   notes.txt                                      │
│                                                     │
│                                                     │
└────────────────────────────────────────────────────┘
 Press Enter to select, Backspace for parent, q to quit
```

### Components

- **Header**: Shows current directory path
- **File List**: Scrollable list of files and directories
- **Status Bar**: Helpful key hints and messages
- **Icons**: for directories, for files

## Keyboard Controls

### Navigation

| Key         | Action                           |
| ----------- | -------------------------------- |
| `↑`         | Move cursor up                   |
| `↓`         | Move cursor down                 |
| `k`         | Move cursor up (Vim-style)       |
| `j`         | Move cursor down (Vim-style)     |
| `Enter`     | Enter directory or download file |
| `Backspace` | Go to parent directory           |
| `q`         | Quit browser                     |
| `Ctrl+C`    | Force quit                       |

### Mouse Support

Currently not supported. Use keyboard navigation.

## File Operations

### Browsing Directories

1. **Navigate** to a directory using arrow keys
2. **Press Enter** to open the directory
3. **View contents** in the updated list
4. **Press Backspace** to go back

Example navigation:

```
/                    (root)
├── documents/       (Enter to open)
│   ├── work/        (Enter to open)
│   │   └── report.pdf
│   └── personal/
└── photos/
```

### Downloading Files

1. **Navigate** to the file you want
2. **Press Enter** on the file
3. **File downloads** to your current directory
4. **Status message** confirms download

```
Downloading: report.pdf...
Downloaded: report.pdf (1.2 MB)
```

### Download Location

Files download to the directory where you ran `orb connect`:

```bash
cd ~/Downloads
orb connect --session <ID> --passcode <CODE>
# Files download to ~/Downloads
```

### Multiple Downloads

Download multiple files by:

1. Download first file
2. Navigate to next file
3. Download next file
4. Repeat as needed

**Note:** Bulk download not currently supported.

## Status Messages

### Success Messages

```
[OK] Downloaded: file.txt
[OK] Directory loaded
[OK] Ready
```

### Error Messages

```
[ERROR] Failed to download: permission denied
[ERROR] Cannot read directory
[ERROR] Connection lost
```

### Loading States

```
Loading directory...
Downloading file...
Connecting...
```

## Advanced Features

### Filtering (Future)

Currently not implemented. Future versions may support:

- Search by filename
- Filter by extension
- Sort by size/date

### Bookmarks (Future)

Currently not implemented. Future versions may support:

- Save frequently accessed paths
- Quick navigation to bookmarks

### Preview (Future)

Currently not implemented. Future versions may support:

- Text file preview
- Image preview
- File size display

## File Types

### Directories

Displayed with icon:

```
 documents/
 photos/
 projects/
```

Press Enter to navigate into directory.

### Regular Files

Displayed with icon:

```
 README.md
 report.pdf
 data.json
```

Press Enter to download file.

### Hidden Files

Hidden files (starting with `.`) are shown:

```
 .gitignore
 .env
 .config/
```

### Symlinks

Symlinks appear as their target type:

- Symlink to directory: Shows as directory
- Symlink to file: Shows as file
- Broken symlink: May show error

## Performance

### Large Directories

For directories with many files:

- Listing may take a moment
- Scrolling remains smooth
- All files shown in list

### Large Files

When downloading large files:

- Download progress not shown (yet)
- Browser may appear frozen
- Wait for completion message

**Tip:** For very large files, consider using a download manager or resumable protocol in future versions.

### Network Latency

On slow connections:

- Directory listings take longer
- File downloads are slower
- Connection may timeout

**Tip:** Keep connection active, avoid network interruptions.

## Troubleshooting

### Browser Won't Open

```
Error: failed to load directory
```

**Causes:**

- Connection failed
- Handshake incomplete
- Permission denied

**Solutions:**

- Check connection logs
- Verify credentials
- Ensure share is active

### Cannot Download File

```
Error: download failed
```

**Causes:**

- File not readable
- Permission denied
- Network interruption
- Disk full

**Solutions:**

- Check file permissions on sharer
- Verify disk space
- Test network connection
- Try smaller file first

### Frozen Interface

**Causes:**

- Large file download in progress
- Network timeout
- Terminal too small

**Solutions:**

- Wait for operation to complete
- Press Ctrl+C to cancel
- Resize terminal window
- Check network connectivity

### Display Issues

**Causes:**

- Terminal doesn't support colors
- Terminal too small
- Unicode not supported

**Solutions:**

- Use modern terminal (iTerm2, Alacritty, Windows Terminal)
- Resize to at least 80x24
- Enable UTF-8 support

## Tips and Tricks

### 1. Terminal Size

For best experience:

- Minimum: 80 columns × 24 rows
- Recommended: 120 columns × 30 rows
- Wide terminals show more files

### 2. Download Organization

```bash
# Create session-specific download directory
mkdir -p ~/Downloads/session-$(date +%Y%m%d)
cd ~/Downloads/session-$(date +%Y%m%d)
orb connect --session <ID> --passcode <CODE>
```

### 3. Quick Navigation

```bash
# Navigate efficiently
# Use ↓↓↓ or jjj to move down quickly
# Press Enter to open
# Backspace to go back
```

### 4. Batch Operations

```bash
# For bulk downloads, script it (future feature)
# Current: manual download each file
```

## Security Considerations

### File Verification

After download, verify:

```bash
# Check file size
ls -lh downloaded-file.pdf

# Verify with checksum (if available)
sha256sum downloaded-file.pdf
```

### Safe Downloads

- Download to isolated directory
- Scan for malware if needed
- Verify file type
- Check file permissions

### Privacy

- File names visible in browser
- Directory structure revealed
- Consider data sensitivity

## Keyboard Shortcuts Reference

Quick reference card:

```
Navigation:
  ↑/k     Move up
  ↓/j     Move down
  Enter   Select/Open/Download
  Bksp    Parent directory

Actions:
  q       Quit browser
  Ctrl+C  Force quit

(More shortcuts in future versions)
```

## Examples

### Navigate to Nested File

```
1. Start at /
2. ↓ to "documents/"
3. Enter
4. ↓ to "projects/"
5. Enter
6. ↓ to "report.pdf"
7. Enter to download
```

### Download Multiple Files

```
1. Navigate to directory
2. ↓ to first file
3. Enter (downloads)
4. ↓ to next file
5. Enter (downloads)
6. Repeat
```

### Explore Directory Structure

```
1. Navigate directories with Enter
2. Backspace to go up
3. Explore branches
4. Backspace back to root
5. Explore other branches
```

## Comparison to Alternatives

### vs. CLI Download

**TUI Browser:**

- Interactive browsing
- Visual directory structure
- Easy navigation
- No automation

**CLI Download:**

- Scriptable
- Batch operations
- No browsing (yet)
- Manual path specification

### vs. GUI Client

**TUI Browser:**

- Works over SSH
- No GUI required
- Lightweight
- Limited features

**GUI Client (future):**

- Mouse support
- Drag and drop
- Requires desktop
- Not implemented

## Future Enhancements

Planned features:

- Search and filter
- Bulk download
- Upload support (if write access added)
- File preview
- Progress bars
- Sorting options
- Bookmarks
- Mouse support
- Copy/paste paths

## Next Steps

- Learn about [Command Reference](commands.md)
- Read [Troubleshooting](troubleshooting.md) guide
- Explore [Connection](connecting.md) details
