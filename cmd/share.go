package cmd

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Zayan-Mohamed/orb/internal/filesystem"
	"github.com/Zayan-Mohamed/orb/internal/tunnel"
	"github.com/Zayan-Mohamed/orb/pkg/protocol"
	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share <path>",
	Short: "Share a local directory",
	Long:  `Share a local directory over an encrypted tunnel. Creates a session ID and passcode.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runShare,
}

var (
	relayURL string
	readOnly bool
)

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().StringVar(&relayURL, "relay", "http://localhost:8080", "Relay server URL")
	shareCmd.Flags().BoolVar(&readOnly, "readonly", false, "Share folder in read-only mode")
}

func runShare(cmd *cobra.Command, args []string) error {
	sharePath := args[0]

	// Validate path exists
	absPath, err := filepath.Abs(sharePath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("path does not exist: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path must be a directory")
	}

	// Create session with relay
	sessionID, passcode, err := createSession(relayURL, absPath)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Display session info
	fmt.Printf("\n")
	fmt.Printf("╔════════════════════════════════════════╗\n")
	fmt.Printf("║     Orb - Secure Folder Sharing       ║\n")
	fmt.Printf("╚════════════════════════════════════════╝\n")
	fmt.Printf("\n")
	fmt.Printf("  Session:  %s\n", sessionID)
	fmt.Printf("  Passcode: %s\n", passcode)
	fmt.Printf("\n")
	fmt.Printf("Share these credentials with the receiver.\n")
	fmt.Printf("Waiting for connection...\n")
	fmt.Printf("\n")

	// Initialize secure filesystem
	secureFS, err := filesystem.NewSecureFilesystem(absPath, readOnly)
	if err != nil {
		return fmt.Errorf("failed to initialize filesystem: %w", err)
	}

	// Connect to relay and establish tunnel
	// Sharer is the responder (waits for connector to initiate handshake)
	tun, err := tunnel.NewTunnel(relayURL, sessionID, passcode, false)
	if err != nil {
		return fmt.Errorf("failed to establish tunnel: %w", err)
	}
	defer tun.Close()

	fmt.Printf("✓ Connected! Tunnel established.\n")
	if readOnly {
		fmt.Printf("  Mode: Read-only\n")
	} else {
		fmt.Printf("  Mode: Read-write\n")
	}
	fmt.Printf("\n")
	fmt.Printf("Press Ctrl+C to stop sharing.\n")
	fmt.Printf("\n")

	// Handle requests
	return handleShareRequests(tun, secureFS)
}

func handleShareRequests(tun *tunnel.Tunnel, fs *filesystem.SecureFilesystem) error {
	for {
		// Receive request
		frame, err := tun.ReceiveFrame()
		if err != nil {
			if tun.IsClosed() {
				return nil
			}
			log.Printf("Error receiving frame: %v", err)
			continue
		}

		// Handle request
		response := processRequest(frame, fs)

		// Send response
		if err := tun.SendFrame(response); err != nil {
			log.Printf("Error sending response: %v", err)
			continue
		}
	}
}

func processRequest(frame *protocol.Frame, fs *filesystem.SecureFilesystem) *protocol.Frame {
	switch frame.Type {
	case protocol.FrameTypePing:
		return &protocol.Frame{
			Type:    protocol.FrameTypePong,
			Payload: []byte{},
		}

	case protocol.FrameTypeList:
		var req protocol.ListRequest
		if err := gob.NewDecoder(bytes.NewReader(frame.Payload)).Decode(&req); err != nil {
			return errorFrame(protocol.ErrCodeUnknown, err.Error())
		}

		resp, err := fs.List(req.Path)
		if err != nil {
			return errorFrame(protocol.ErrCodeIO, err.Error())
		}

		return responseFrame(resp)

	case protocol.FrameTypeStat:
		var req protocol.StatRequest
		if err := gob.NewDecoder(bytes.NewReader(frame.Payload)).Decode(&req); err != nil {
			return errorFrame(protocol.ErrCodeUnknown, err.Error())
		}

		resp, err := fs.Stat(req.Path)
		if err != nil {
			return errorFrame(protocol.ErrCodeNotFound, err.Error())
		}

		return responseFrame(resp)

	case protocol.FrameTypeRead:
		var req protocol.ReadRequest
		if err := gob.NewDecoder(bytes.NewReader(frame.Payload)).Decode(&req); err != nil {
			return errorFrame(protocol.ErrCodeUnknown, err.Error())
		}

		resp, err := fs.Read(req.Path, req.Offset, req.Length)
		if err != nil {
			return errorFrame(protocol.ErrCodeIO, err.Error())
		}

		return responseFrame(resp)

	case protocol.FrameTypeWrite:
		var req protocol.WriteRequest
		if err := gob.NewDecoder(bytes.NewReader(frame.Payload)).Decode(&req); err != nil {
			return errorFrame(protocol.ErrCodeUnknown, err.Error())
		}

		resp, err := fs.Write(req.Path, req.Offset, req.Data)
		if err != nil {
			return errorFrame(protocol.ErrCodePermission, err.Error())
		}

		return responseFrame(resp)

	case protocol.FrameTypeDelete:
		var req protocol.DeleteRequest
		if err := gob.NewDecoder(bytes.NewReader(frame.Payload)).Decode(&req); err != nil {
			return errorFrame(protocol.ErrCodeUnknown, err.Error())
		}

		if err := fs.Delete(req.Path); err != nil {
			return errorFrame(protocol.ErrCodePermission, err.Error())
		}

		return responseFrame(&protocol.WriteResponse{BytesWritten: 0})

	case protocol.FrameTypeRename:
		var req protocol.RenameRequest
		if err := gob.NewDecoder(bytes.NewReader(frame.Payload)).Decode(&req); err != nil {
			return errorFrame(protocol.ErrCodeUnknown, err.Error())
		}

		if err := fs.Rename(req.OldPath, req.NewPath); err != nil {
			return errorFrame(protocol.ErrCodePermission, err.Error())
		}

		return responseFrame(&protocol.WriteResponse{BytesWritten: 0})

	case protocol.FrameTypeMkdir:
		var req protocol.MkdirRequest
		if err := gob.NewDecoder(bytes.NewReader(frame.Payload)).Decode(&req); err != nil {
			return errorFrame(protocol.ErrCodeUnknown, err.Error())
		}

		if err := fs.Mkdir(req.Path, req.Perm); err != nil {
			return errorFrame(protocol.ErrCodePermission, err.Error())
		}

		return responseFrame(&protocol.WriteResponse{BytesWritten: 0})

	default:
		return errorFrame(protocol.ErrCodeUnknown, "unknown request type")
	}
}

func responseFrame(data interface{}) *protocol.Frame {
	var buf bytes.Buffer
	_ = gob.NewEncoder(&buf).Encode(data)

	return &protocol.Frame{
		Type:    protocol.FrameTypeResponse,
		Payload: buf.Bytes(),
	}
}

func errorFrame(code uint32, message string) *protocol.Frame {
	errResp := protocol.ErrorResponse{
		Code:    code,
		Message: message,
	}

	var buf bytes.Buffer
	_ = gob.NewEncoder(&buf).Encode(errResp)

	return &protocol.Frame{
		Type:    protocol.FrameTypeError,
		Payload: buf.Bytes(),
	}
}
