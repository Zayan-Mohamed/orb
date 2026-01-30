package tui

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Zayan-Mohamed/orb/internal/tunnel"
	"github.com/Zayan-Mohamed/orb/pkg/protocol"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Download progress messages
type downloadProgressMsg struct {
	downloaded int64
	speed      int64
}

type downloadCompleteMsg struct {
	filename string
	size     int64
}

type downloadErrorMsg struct {
	error string
}

type downloadCancelMsg struct{}

type downloadResetMsg struct{}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 0)

	progressStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46")).
			Padding(0, 1)

	progressBarStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Background(lipgloss.Color("240")).
				Width(50)

	progressFilledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("46")).
				Background(lipgloss.Color("46"))
)

type downloadState struct {
	filename      string
	totalSize     int64
	downloaded    int64
	chunkSize     int64
	isDownloading bool
	cancelled     bool
	progress      float64
	speed         int64 // bytes per second
	startTime     int64 // Unix timestamp
}

type fileItem struct {
	name  string
	size  int64
	isDir bool
}

func (i fileItem) Title() string {
	if i.isDir {
		return "📁 " + i.name
	}
	return "📄 " + i.name
}

func (i fileItem) Description() string {
	if i.isDir {
		return "<DIR>"
	}
	return formatSize(i.size)
}

func (i fileItem) FilterValue() string {
	return i.name
}

type model struct {
	tunnel      *tunnel.Tunnel
	currentPath string
	list        list.Model
	error       string
	download    downloadState // NEW: Add download state
}

func newModel(tun *tunnel.Tunnel) model {
	items := []list.Item{}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Orb File Browser"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle

	return model{
		tunnel:      tun,
		currentPath: "/",
		list:        l,
		download:    downloadState{}, // Initialize download state
	}
}

func (m model) Init() tea.Cmd {
	return m.loadDirectory()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle download progress messages first (highest priority)
	switch msg := msg.(type) {
	case downloadProgressMsg:
		if m.download.isDownloading && !m.download.cancelled {
			m.download.downloaded = msg.downloaded
			m.download.speed = msg.speed
			m.download.progress = float64(msg.downloaded) / float64(m.download.totalSize) * 100
			return m, nil
		}

	case downloadCompleteMsg:
		m.download.isDownloading = false
		m.download.progress = 100
		// Reset after 2 seconds
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return downloadResetMsg{}
		})

	case downloadErrorMsg:
		m.download.isDownloading = false
		m.error = msg.error
		return m, nil

	case downloadCancelMsg:
		// Reset download state
		m.download = downloadState{}
		return m, m.loadDirectory()

	case downloadResetMsg:
		// Reset download state
		m.download = downloadState{}
		return m, m.loadDirectory()
	}

	// Handle key messages with download cancellation
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		// ESC key cancels downloads
		if key.Matches(msg, key.NewBinding(key.WithKeys("escape"))) {
			if m.download.isDownloading {
				m.download.cancelled = true
				m.download.isDownloading = false
				return m, nil
			}
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q"))):
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.download.isDownloading {
				return m, nil // Ignore input during download
			}

			selected := m.list.SelectedItem()
			if selected != nil {
				item := selected.(fileItem)
				if item.isDir {
					if item.name == ".." {
						m.currentPath = filepath.Dir(m.currentPath)
					} else {
						m.currentPath = filepath.Join(m.currentPath, item.name)
					}
					return m, m.loadDirectory()
				} else {
					return m, m.initiateDownload(item.name, item.size)
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
			if m.download.isDownloading {
				return m, nil
			}
			if m.currentPath != "/" {
				m.currentPath = filepath.Dir(m.currentPath)
				return m, m.loadDirectory()
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("d"))):
			if m.download.isDownloading {
				return m, nil
			}

			selected := m.list.SelectedItem()
			if selected != nil {
				item := selected.(fileItem)
				if !item.isDir {
					return m, m.initiateDownload(item.name, item.size)
				}
			}
		}

	case []list.Item:
		if !m.download.isDownloading {
			m.list.SetItems(msg)
			m.error = ""
		}
		return m, nil

	case error:
		if !m.download.isDownloading {
			m.error = msg.Error()
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	// Show progress overlay during download
	if m.download.isDownloading {
		b.WriteString(m.renderDownloadProgress())
		return b.String()
	}

	// Title
	b.WriteString(m.list.View())
	b.WriteString("\n")

	// Current path
	b.WriteString(statusStyle.Render("Path: " + m.currentPath))
	b.WriteString("\n")

	// Error message
	if m.error != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("Error: " + m.error))
		b.WriteString("\n")
	}

	// Help
	helpText := "Enter: open/download • d: download • backspace: parent dir"
	if m.download.isDownloading {
		helpText = "ESC: cancel download"
	}
	helpText += " • q: quit"
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

func (m model) renderDownloadProgress() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Downloading File"))
	b.WriteString("\n\n")

	// Filename
	b.WriteString(progressStyle.Render("File: " + m.download.filename))
	b.WriteString("\n")

	// Progress bar
	barWidth := 50
	filled := int(float64(barWidth) * m.download.progress / 100)
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	progressBar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	b.WriteString(progressBarStyle.Render(progressBar))
	b.WriteString("\n\n")

	// Progress info
	progressText := fmt.Sprintf("%.1f%%", m.download.progress)
	sizeText := fmt.Sprintf("%s / %s",
		formatSize(m.download.downloaded),
		formatSize(m.download.totalSize))

	b.WriteString(progressStyle.Render(progressText + "  " + sizeText))
	b.WriteString("\n")

	// Speed
	if m.download.speed > 0 {
		speedText := fmt.Sprintf("Speed: %s/s", formatSize(m.download.speed))
		b.WriteString(statusStyle.Render(speedText))
		b.WriteString("\n")
	}

	// Cancel hint
	b.WriteString(helpStyle.Render("Press ESC to cancel"))

	return b.String()
}

func (m model) loadDirectory() tea.Cmd {
	return func() tea.Msg {
		// Send list request
		req := protocol.ListRequest{
			Path: m.currentPath,
		}

		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(req); err != nil {
			return err
		}

		frame := &protocol.Frame{
			Type:    protocol.FrameTypeList,
			Payload: buf.Bytes(),
		}

		if err := m.tunnel.SendFrame(frame); err != nil {
			return err
		}

		// Receive response
		respFrame, err := m.tunnel.ReceiveFrame()
		if err != nil {
			return err
		}

		if respFrame.Type == protocol.FrameTypeError {
			var errResp protocol.ErrorResponse
			_ = gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&errResp)
			return fmt.Errorf("%s", errResp.Message)
		}

		if respFrame.Type != protocol.FrameTypeResponse {
			return fmt.Errorf("unexpected frame type: %d", respFrame.Type)
		}

		var resp protocol.ListResponse
		if err := gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&resp); err != nil {
			return err
		}

		// Convert to list items
		items := []list.Item{}

		// Add parent directory entry if not at root
		if m.currentPath != "/" {
			items = append(items, fileItem{
				name:  "..",
				isDir: true,
			})
		}

		for _, file := range resp.Files {
			items = append(items, fileItem{
				name:  file.Name,
				size:  file.Size,
				isDir: file.IsDir,
			})
		}

		return items
	}
}

func (m model) initiateDownload(filename string, size int64) tea.Cmd {
	return func() tea.Msg {
		// Initialize download state
		m.download.filename = filename
		m.download.totalSize = size
		m.download.chunkSize = 64 * 1024 // 64KB chunks
		m.download.downloaded = 0
		m.download.isDownloading = true
		m.download.cancelled = false
		m.download.progress = 0
		m.download.startTime = time.Now().Unix()

		remotePath := filepath.Join(m.currentPath, filename)
		localPath := filepath.Join(".", filename)

		// Create local file
		file, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			return downloadErrorMsg{error: err.Error()}
		}
		defer file.Close()

		// Download in chunks
		var totalDownloaded int64
		chunkSize := m.download.chunkSize

		for offset := int64(0); offset < size; offset += chunkSize {
			// Check for cancellation
			if m.download.cancelled {
				os.Remove(localPath) // Clean up partial file
				return downloadCancelMsg{}
			}

			// Calculate chunk size for this iteration
			remaining := size - offset
			currentChunkSize := chunkSize
			if remaining < chunkSize {
				currentChunkSize = remaining
			}

			// Send read request for this chunk
			readReq := protocol.ReadRequest{
				Path:   remotePath,
				Offset: offset,
				Length: currentChunkSize,
			}

			var buf bytes.Buffer
			if err := gob.NewEncoder(&buf).Encode(readReq); err != nil {
				return downloadErrorMsg{error: err.Error()}
			}

			frame := &protocol.Frame{
				Type:    protocol.FrameTypeRead,
				Payload: buf.Bytes(),
			}

			if err := m.tunnel.SendFrame(frame); err != nil {
				return downloadErrorMsg{error: err.Error()}
			}

			// Receive chunk response
			respFrame, err := m.tunnel.ReceiveFrame()
			if err != nil {
				return downloadErrorMsg{error: err.Error()}
			}

			if respFrame.Type == protocol.FrameTypeError {
				var errResp protocol.ErrorResponse
				_ = gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&errResp)
				return downloadErrorMsg{error: errResp.Message}
			}

			if respFrame.Type != protocol.FrameTypeResponse {
				return downloadErrorMsg{error: fmt.Sprintf("unexpected frame type: %d", respFrame.Type)}
			}

			var readResp protocol.ReadResponse
			if err := gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&readResp); err != nil {
				return downloadErrorMsg{error: err.Error()}
			}

			// Write chunk to file
			if _, err := file.WriteAt(readResp.Data, offset); err != nil {
				return downloadErrorMsg{error: err.Error()}
			}

			totalDownloaded += int64(len(readResp.Data))

			// Calculate speed (bytes per second)
			elapsed := time.Now().Unix() - m.download.startTime
			var speed int64
			if elapsed > 0 {
				speed = totalDownloaded / elapsed
			}

			// Send progress update
			return downloadProgressMsg{
				downloaded: totalDownloaded,
				speed:      speed,
			}
		}

		// Download complete
		return downloadCompleteMsg{
			filename: filename,
			size:     totalDownloaded,
		}
	}
}

func (m model) downloadFile(filename string) tea.Cmd {
	// Legacy method - now redirects to initiateDownload
	// We need to get file size first
	return func() tea.Msg {
		remotePath := filepath.Join(m.currentPath, filename)

		// Get file info first
		statReq := protocol.StatRequest{
			Path: remotePath,
		}

		var buf bytes.Buffer
		_ = gob.NewEncoder(&buf).Encode(statReq)

		frame := &protocol.Frame{
			Type:    protocol.FrameTypeStat,
			Payload: buf.Bytes(),
		}

		if err := m.tunnel.SendFrame(frame); err != nil {
			return err
		}

		respFrame, err := m.tunnel.ReceiveFrame()
		if err != nil {
			return err
		}

		if respFrame.Type == protocol.FrameTypeError {
			var errResp protocol.ErrorResponse
			_ = gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&errResp)
			return fmt.Errorf("%s", errResp.Message)
		}

		if respFrame.Type != protocol.FrameTypeResponse {
			return fmt.Errorf("unexpected frame type: %d", respFrame.Type)
		}

		var statResp protocol.StatResponse
		_ = gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&statResp)

		// Now initiate the actual download
		cmd := m.initiateDownload(filename, statResp.Info.Size)
		return cmd()
	}
}

// StartFileBrowser starts the TUI file browser
func StartFileBrowser(tun *tunnel.Tunnel) error {
	m := newModel(tun)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
