package main

import (
	"fmt"
	"os"

	"github.com/csawai/gitx/internal/config"
	"github.com/csawai/gitx/internal/keychain"
	"github.com/csawai/gitx/internal/ssh"
	"github.com/spf13/cobra"
)

var removeIdentityCmd = &cobra.Command{
	Use:   "remove identity [alias]",
	Short: "Remove an identity",
	Long:  "Remove an identity and clean up associated SSH config entries and keychain secrets.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := removeIdentity(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func removeIdentity(alias string) error {
	// Load identity to get SSH host alias
	identity, err := config.FindIdentityByAlias(alias)
	if err != nil {
		return err
	}

	// Remove SSH config entry if exists
	if identity.SSHHostAlias != "" {
		if err := ssh.RemoveSSHConfigEntry(identity.SSHHostAlias); err != nil {
			// Log but don't fail
			fmt.Fprintf(os.Stderr, "Warning: failed to remove SSH config entry: %v\n", err)
		}
	}

	// Remove keychain secrets
	if err := keychain.DeleteAllSecrets(alias); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to remove keychain secrets: %v\n", err)
	}

	// Remove from config
	if err := config.RemoveIdentity(alias); err != nil {
		return err
	}

	fmt.Printf("✓ Identity '%s' removed\n", alias)
	return nil
}
