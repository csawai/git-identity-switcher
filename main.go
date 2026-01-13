package main

import (
	"fmt"
	"os"

	"github.com/csawai/git-identity-switcher/internal/ui"
	"github.com/spf13/cobra"
)

func showBanner() {
	fmt.Println(ui.Banner())
	fmt.Println(ui.Subtitle("Git Identity Switcher - Never push to the wrong account again"))
	fmt.Println()
}

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "gitx",
	Short: "Git Identity Switcher - Manage multiple GitHub identities safely",
	Long: `gitx is a CLI tool for managing multiple GitHub identities with per-repo binding.
It helps developers safely switch between work, personal, and client accounts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show banner and help if no command provided
		showBanner()
		cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show gitx version",
	Run: func(cmd *cobra.Command, args []string) {
		showBanner()
		fmt.Println(ui.InfoBox.Render(fmt.Sprintf("Version: %s", version)))
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
	// show_key.go registers showKeyCmd and copyKeyCmd
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// Use styled error box
		errorMsg := ui.ErrorBox.Render(fmt.Sprintf("❌ Error: %v", err))
		fmt.Fprintf(os.Stderr, "%s\n", errorMsg)
		os.Exit(1)
	}
}

