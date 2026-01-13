package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "gitx",
	Short: "Git Identity Switcher - Manage multiple GitHub identities safely",
	Long: `gitx is a CLI tool for managing multiple GitHub identities with per-repo binding.
It helps developers safely switch between work, personal, and client accounts.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show gitx version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gitx version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(addIdentityCmd)
	rootCmd.AddCommand(listIdentitiesCmd)
	rootCmd.AddCommand(bindCmd)
	rootCmd.AddCommand(unbindCmd)
	rootCmd.AddCommand(removeIdentityCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

