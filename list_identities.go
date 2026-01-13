package main

import (
	"fmt"
	"os"

	"github.com/csawai/gitx/internal/config"
	"github.com/spf13/cobra"
)

var listIdentitiesCmd = &cobra.Command{
	Use:   "list identities",
	Short: "List all stored identities",
	Long:  "Display all configured identities (no secrets shown).",
	Run: func(cmd *cobra.Command, args []string) {
		if err := listIdentities(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func listIdentities() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if len(cfg.Identities) == 0 {
		fmt.Println("No identities configured.")
		return nil
	}

	fmt.Println("Configured Identities:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, id := range cfg.Identities {
		fmt.Printf("Alias:      %s\n", id.Alias)
		fmt.Printf("Name:       %s\n", id.Name)
		fmt.Printf("Email:      %s\n", id.Email)
		fmt.Printf("GitHub:     %s\n", id.GitHubUser)
		fmt.Printf("Auth:       %s\n", id.AuthMethod)
		if id.SSHKeyPath != "" {
			fmt.Printf("SSH Key:    %s\n", id.SSHKeyPath)
		}
		fmt.Println()
	}

	return nil
}

