package tui

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Zayan-Mohamed/orb/internal/tunnel"
	"github.com/Zayan-Mohamed/orb/pkg/protocol"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
)

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
	downloading bool
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
	}
}

func (m model) Init() tea.Cmd {
	return m.loadDirectory()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q"))):
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			selected := m.list.SelectedItem()
			if selected != nil {
				item := selected.(fileItem)
				if item.isDir {
					// Navigate into directory
					if item.name == ".." {
						m.currentPath = filepath.Dir(m.currentPath)
					} else {
						m.currentPath = filepath.Join(m.currentPath, item.name)
					}
					return m, m.loadDirectory()
				} else {
					// Download file
					return m, m.downloadFile(item.name)
				}
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
			// Go up one directory
			if m.currentPath != "/" {
				m.currentPath = filepath.Dir(m.currentPath)
				return m, m.loadDirectory()
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("d"))):
			// Download selected file
			selected := m.list.SelectedItem()
			if selected != nil {
				item := selected.(fileItem)
				if !item.isDir {
					return m, m.downloadFile(item.name)
				}
			}
		}

	case []list.Item:
		m.list.SetItems(msg)
		m.error = ""
		return m, nil

	case error:
		m.error = msg.Error()
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

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
	b.WriteString(helpStyle.Render("Enter: open/download • d: download • backspace: parent dir • q: quit"))

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
			gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&errResp)
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

func (m model) downloadFile(filename string) tea.Cmd {
	return func() tea.Msg {
		remotePath := filepath.Join(m.currentPath, filename)

		// Get file info first
		statReq := protocol.StatRequest{
			Path: remotePath,
		}

		var buf bytes.Buffer
		gob.NewEncoder(&buf).Encode(statReq)

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
			gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&errResp)
			return fmt.Errorf("%s", errResp.Message)
		}

		if respFrame.Type != protocol.FrameTypeResponse {
			return fmt.Errorf("unexpected frame type: %d", respFrame.Type)
		}

		var statResp protocol.StatResponse
		gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&statResp)

		// Read file content
		readReq := protocol.ReadRequest{
			Path:   remotePath,
			Offset: 0,
			Length: statResp.Info.Size,
		}

		buf.Reset()
		gob.NewEncoder(&buf).Encode(readReq)

		frame = &protocol.Frame{
			Type:    protocol.FrameTypeRead,
			Payload: buf.Bytes(),
		}

		if err := m.tunnel.SendFrame(frame); err != nil {
			return err
		}

		respFrame, err = m.tunnel.ReceiveFrame()
		if err != nil {
			return err
		}

		if respFrame.Type == protocol.FrameTypeError {
			var errResp protocol.ErrorResponse
			gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&errResp)
			return fmt.Errorf("%s", errResp.Message)
		}

		if respFrame.Type != protocol.FrameTypeResponse {
			return fmt.Errorf("unexpected frame type: %d", respFrame.Type)
		}

		var readResp protocol.ReadResponse
		gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&readResp)

		// Save to local file
		localPath := filepath.Join(".", filename)
		if err := os.WriteFile(localPath, readResp.Data, 0644); err != nil {
			return err
		}

		// Reload directory to refresh
		return m.loadDirectory()()
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

// UploadFile uploads a file to the remote filesystem
func (m model) uploadFile(localPath string) tea.Cmd {
	return func() tea.Msg {
		// Read local file
		data, err := os.ReadFile(localPath)
		if err != nil {
			return err
		}

		// Upload to remote
		filename := filepath.Base(localPath)
		remotePath := filepath.Join(m.currentPath, filename)

		writeReq := protocol.WriteRequest{
			Path:   remotePath,
			Offset: 0,
			Data:   data,
		}

		var buf bytes.Buffer
		gob.NewEncoder(&buf).Encode(writeReq)

		frame := &protocol.Frame{
			Type:    protocol.FrameTypeWrite,
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
			gob.NewDecoder(bytes.NewReader(respFrame.Payload)).Decode(&errResp)
			return fmt.Errorf("%s", errResp.Message)
		}

		// Reload directory
		return m.loadDirectory()()
	}
}

// Helper to read input from user
func readInput(prompt string) (string, error) {
	fmt.Print(prompt)
	var input string
	_, err := fmt.Scanln(&input)
	return input, err
}
