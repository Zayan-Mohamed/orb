# Creating GIFs for Documentation

This guide helps you create high-quality GIF screen recordings for the Orb documentation.

## Required Tools

### Option 1: Using asciinema + agg (Recommended for terminal recordings)

```bash
# Install asciinema
# Ubuntu/Debian
sudo apt-get install asciinema

# macOS
brew install asciinema

# Install agg (converts asciinema to GIF)
cargo install --git https://github.com/asciinema/agg
# OR download binary from: https://github.com/asciinema/agg/releases
```

### Option 2: Using ttyd + peek/gifski

```bash
# Ubuntu/Debian
sudo apt-get install peek

# macOS
brew install gifski
```

### Option 3: Using FFmpeg

```bash
# Ubuntu/Debian
sudo apt-get install ffmpeg

# macOS
brew install ffmpeg
```

## Recording Workflow

### Method 1: asciinema + agg (Best Quality)

This method produces the cleanest terminal recordings.

#### 1. Record Terminal Session

```bash
# Record the session
asciinema rec demo.cast

# Now perform your actions:
# - orb share ./myfolder
# - Wait a moment
# - Press Ctrl+D to stop recording
```

#### 2. Convert to GIF

```bash
# Basic conversion
agg demo.cast demo.gif

# High quality with custom settings
agg --cols 80 --rows 24 --speed 1.5 demo.cast demo.gif

# Recommended settings for documentation
agg --cols 100 --rows 30 --speed 2.0 --font-size 16 demo.cast demo.gif
```

#### 3. Optimize GIF Size

```bash
# Install gifsicle for optimization
sudo apt-get install gifsicle  # Ubuntu/Debian
brew install gifsicle          # macOS

# Optimize
gifsicle -O3 --colors 256 demo.gif -o demo-optimized.gif
```

### Method 2: Record Video + Convert with FFmpeg

#### 1. Record Your Screen

- **Linux**: Use `Kazam`, `SimpleScreenRecorder`, or `OBS Studio`
- **macOS**: Use `QuickTime Player` (Cmd+Shift+5) or `OBS Studio`
- **Windows**: Use `OBS Studio` or built-in Game Bar (Win+G)

Tips for recording:

- Use a resolution around 1280x720 for best results
- Keep recordings under 30 seconds
- Use a clean terminal theme (dark background preferred)
- Increase font size for better readability

#### 2. Convert Video to GIF

```bash
# Generate color palette (improves quality)
ffmpeg -i demo.mp4 -vf "fps=10,scale=1000:-1:flags=lanczos,palettegen" palette.png

# Create GIF using the palette
ffmpeg -i demo.mp4 -i palette.png -filter_complex \
  "fps=10,scale=1000:-1:flags=lanczos[x];[x][1:v]paletteuse" demo.gif

# Remove palette file
rm palette.png
```

#### 3. Optimize GIF

```bash
gifsicle -O3 --colors 256 --lossy=80 demo.gif -o demo-optimized.gif
```

### Method 3: Using Peek (Linux Only - Quick & Easy)

```bash
# Install Peek
sudo apt-get install peek

# Run Peek
peek

# 1. Position the recording window over your terminal
# 2. Click Record
# 3. Perform your actions
# 4. Click Stop
# 5. Save as GIF
```

## GIF Specifications for Orb Documentation

### File Naming Convention

```
orb-<action>-demo.gif
```

Examples:

- `orb-share-demo.gif` - Sharing a folder
- `orb-connect-demo.gif` - Connecting to a share
- `orb-complete-demo.gif` - Complete workflow

### Recommended Settings

- **Resolution**: 1000-1200px wide
- **Frame Rate**: 10-15 FPS
- **Duration**: 10-30 seconds
- **File Size**: < 5MB (optimize if larger)
- **Colors**: 256 colors maximum
- **Terminal Size**: 80-100 columns, 24-30 rows

### Content Guidelines

1. **Share Demo** (`orb-share-demo.gif`):

   - Start with empty terminal
   - Run `orb share ./myfolder`
   - Show session ID and passcode output
   - Wait 2-3 seconds showing "Waiting for connection..."
   - End recording

2. **Connect Demo** (`orb-connect-demo.gif`):

   - Start with empty terminal
   - Run `orb connect <session-id>`
   - Enter passcode when prompted
   - Show TUI file browser loading
   - Navigate through a few files
   - Press `q` to quit
   - End recording

3. **Complete Demo** (`orb-complete-demo.gif`):
   - Split screen showing both terminals
   - Left: Share folder
   - Right: Connect and browse
   - Show file download
   - End with both terminals

## Online Tools (No Installation)

If you don't want to install software, you can use:

1. **ScreenToGif** (Windows): https://www.screentogif.com/
2. **GIPHY Capture** (macOS): https://giphy.com/apps/giphycapture
3. **Recordit** (Cross-platform): https://recordit.co/

## Tips for High-Quality GIFs

1. **Use a Clean Terminal Theme**

   ```bash
   # Recommended: Use a simple, high-contrast theme
   # - Dark background
   # - Bright text
   # - Minimal decorations
   ```

2. **Increase Terminal Font Size**

   ```bash
   # Make text larger for better readability
   # In most terminals: Ctrl+Plus
   ```

3. **Keep It Short**

   - 15-20 seconds is ideal
   - Focus on one action per GIF

4. **Add Pauses**

   - Pause 1-2 seconds at important moments
   - Let viewers read output

5. **Optimize File Size**
   ```bash
   # If GIF is > 5MB, reduce colors or frame rate
   gifsicle -O3 --colors 128 --lossy=100 input.gif -o output.gif
   ```

## Placing GIFs in Documentation

Once you have your GIFs:

1. Place them in `docs/assets/images/`
2. Update references in `README.md`
3. Ensure file names match the placeholders

Current placeholders in README:

- `docs/assets/images/orb-share-demo.gif`
- `docs/assets/images/orb-connect-demo.gif`
- `docs/assets/images/orb-complete-demo.gif`

## Testing Your GIFs

Before committing, verify:

- ✅ File size < 5MB
- ✅ Resolution clear and readable
- ✅ Frame rate smooth (10+ FPS)
- ✅ Duration appropriate (10-30 seconds)
- ✅ Colors look good (256 colors max)

## Quick Reference Commands

```bash
# Record with asciinema
asciinema rec demo.cast

# Convert with agg
agg --cols 100 --rows 30 --speed 2.0 demo.cast demo.gif

# Optimize with gifsicle
gifsicle -O3 --colors 256 demo.gif -o demo-optimized.gif

# Check file size
ls -lh demo-optimized.gif
```

## Need Help?

If you have questions about creating GIFs for documentation:

- Check existing GIF examples in `docs/assets/images/`
- Open an issue on GitHub
- Refer to tool-specific documentation

## Alternative: Video Files

If GIFs are too large or difficult to create, you can also:

- Upload videos to YouTube/Vimeo
- Embed video links in documentation
- Use animated PNGs (APNG) as an alternative to GIF
