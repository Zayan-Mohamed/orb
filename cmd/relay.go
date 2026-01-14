package cmd

import (
	"fmt"
	"log"

	"github.com/Zayan-Mohamed/orb/internal/relay"
	"github.com/spf13/cobra"
)

var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "Start the relay server",
	Long:  `Start the relay server that forwards encrypted connections between peers.`,
	RunE:  runRelay,
}

var (
	listenAddr string
)

func init() {
	rootCmd.AddCommand(relayCmd)
	relayCmd.Flags().StringVar(&listenAddr, "listen", ":8080", "Listen address (e.g., :8080 or 0.0.0.0:8080)")
}

func runRelay(cmd *cobra.Command, args []string) error {
	fmt.Printf("Starting Orb relay server...\n")
	fmt.Printf("Listening on %s\n", listenAddr)
	fmt.Printf("\n")
	fmt.Printf("Security notes:\n")
	fmt.Printf("  • The relay server never sees plaintext data\n")
	fmt.Printf("  • All encryption happens at the edges\n")
	fmt.Printf("  • Sessions expire automatically\n")
	fmt.Printf("\n")

	server := relay.NewRelayServer()
	defer server.Shutdown()

	if err := server.Start(listenAddr); err != nil {
		log.Fatalf("Relay server error: %v", err)
	}

	return nil
}
