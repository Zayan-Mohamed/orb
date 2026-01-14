package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Zayan-Mohamed/orb/internal/tui"
	"github.com/Zayan-Mohamed/orb/internal/tunnel"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect <session-id>",
	Short: "Connect to a shared session",
	Long:  `Connect to a shared folder session using the session ID and passcode.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runConnect,
}

var (
	passcode  string
	mountPath string
	tuiMode   bool
)

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVar(&relayURL, "relay", "http://localhost:8080", "Relay server URL")
	connectCmd.Flags().StringVarP(&passcode, "passcode", "p", "", "Session passcode (will prompt if not provided)")
	connectCmd.Flags().StringVarP(&mountPath, "mount", "m", "", "Mount point (Linux/macOS only)")
	connectCmd.Flags().BoolVar(&tuiMode, "tui", true, "Use TUI file browser")
}

func runConnect(cmd *cobra.Command, args []string) error {
	sessionID := args[0]

	// Prompt for passcode if not provided
	if passcode == "" {
		fmt.Print("Enter passcode: ")
		_, _ = fmt.Scanln(&passcode)
	}

	// Establish tunnel
	fmt.Printf("Connecting to session %s...\n", sessionID)

	// Connector is the initiator (starts the handshake)
	tun, err := tunnel.NewTunnel(relayURL, sessionID, passcode, true)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		if err := tun.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close tunnel: %v\n", err)
		}
	}()

	fmt.Printf("✓ Connected! Tunnel established.\n")

	// Determine mode based on platform and flags
	canMount := runtime.GOOS == "linux" || runtime.GOOS == "darwin"

	if mountPath != "" && canMount {
		// Try FUSE mounting (Linux/macOS only)
		fmt.Printf("Mounting at %s...\n", mountPath)
		return mountFilesystem(tun, mountPath)
	}

	// Use TUI file browser (cross-platform)
	if tuiMode {
		fmt.Printf("Opening file browser...\n")
		fmt.Printf("Press Ctrl+C to disconnect.\n\n")
		return tui.StartFileBrowser(tun)
	}

	return fmt.Errorf("no mode selected (use --tui or --mount)")
}

// mountFilesystem mounts the remote filesystem using FUSE
func mountFilesystem(tun *tunnel.Tunnel, mountPoint string) error {
	// This will be implemented with FUSE support
	return fmt.Errorf("FUSE mounting not yet implemented - use --tui mode")
}
