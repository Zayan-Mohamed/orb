package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information - set during build with -ldflags
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "orb",
	Short: "Orb - Zero-Trust Folder Tunneling Tool",
	Long: `Orb is a secure folder sharing tool that uses end-to-end encryption.
No accounts, no cloud storage, no port forwarding.
All data is encrypted and the relay server is blind.`,
	Version: Version,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Orb version %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Build date: %s\n", BuildDate)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetVersionTemplate(fmt.Sprintf("Orb version %s\nGit commit: %s\nBuild date: %s\n", Version, GitCommit, BuildDate))
	rootCmd.AddCommand(versionCmd)
}
