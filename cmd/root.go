package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "orb",
	Short: "Orb - Zero-Trust Folder Tunneling Tool",
	Long: `Orb is a secure folder sharing tool that uses end-to-end encryption.
No accounts, no cloud storage, no port forwarding.
All data is encrypted and the relay server is blind.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
