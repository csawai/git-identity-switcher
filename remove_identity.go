package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/csawai/git-identity-switcher/internal/config"
	"github.com/csawai/git-identity-switcher/internal/keychain"
	"github.com/csawai/git-identity-switcher/internal/ssh"
	"github.com/csawai/git-identity-switcher/internal/ui"
	"github.com/spf13/cobra"
)

var (
	removeDryRun      bool
	removeForce        bool
	removeDeleteKeys   bool
)

func init() {
	removeIdentityCmd.Flags().BoolVar(&removeDryRun, "dry-run", false, "Show what would be deleted without making changes")
	removeIdentityCmd.Flags().BoolVar(&removeForce, "force", false, "Skip confirmation prompt")
	removeIdentityCmd.Flags().BoolVar(&removeDeleteKeys, "delete-keys", false, "Delete SSH key files")
}

var removeIdentityCmd = &cobra.Command{
	Use:   "remove identity [alias]",
	Short: "Remove an identity",
	Long:  "Remove an identity and clean up associated SSH config entries, keychain secrets, and optionally SSH key files.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := removeIdentity(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func removeIdentity(alias string) error {
	// Load identity to get details
	identity, err := config.FindIdentityByAlias(alias)
	if err != nil {
		return err
	}

	// Show what will be deleted in a styled box
	var items []string
	items = append(items, "• Identity config entry")
	
	if identity.SSHHostAlias != "" {
		items = append(items, fmt.Sprintf("• SSH config entry: %s", identity.SSHHostAlias))
	}
	
	if identity.AuthMethod == "pat" {
		items = append(items, "• Keychain secrets (PAT)")
		items = append(items, "• Git credential helper entries")
	}
	
	if identity.SSHKeyPath != "" {
		items = append(items, fmt.Sprintf("• SSH key files: %s, %s.pub", identity.SSHKeyPath, identity.SSHKeyPath))
	}

	content := fmt.Sprintf("🗑️  Removing identity: %s\n\n", ui.ErrorText.Render(alias))
	content += "This will remove:\n"
	for _, item := range items {
		content += "  " + item + "\n"
	}
	content += "\n" + ui.WarningText.Render("⚠️  Warning:") + " If any repositories are bound to this identity\n"
	content += fmt.Sprintf("   (email: %s), you may need to rebind them.", identity.Email)

	fmt.Println(ui.WarningBox.Render(content))
	fmt.Println()

	// Dry-run mode
	if removeDryRun {
		fmt.Println("[DRY RUN] No changes were made.")
		return nil
	}

	// Confirmation
	if !removeForce {
		fmt.Print("Are you sure you want to remove this identity? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Ask about SSH keys if not specified via flag
	deleteKeys := removeDeleteKeys
	if identity.SSHKeyPath != "" && !removeDeleteKeys && !removeForce {
		fmt.Print("Delete SSH key files? (y/n) [n]: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response == "y" || response == "yes" {
			deleteKeys = true
		}
	}

		// Remove SSH config entry
	if identity.SSHHostAlias != "" {
		if err := ui.SpinnerWithFunc("Removing SSH config entry", func() error {
			return ssh.RemoveSSHConfigEntry(identity.SSHHostAlias)
		}); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", ui.WarningText.Render("⚠️  Warning: failed to remove SSH config entry: "+err.Error()))
		}
	}

	// Remove keychain secrets and git credentials
	if identity.AuthMethod == "pat" {
		// Remove from gitx keychain
		if err := ui.SpinnerWithFunc("Removing keychain secrets", func() error {
			return keychain.DeleteAllSecrets(alias)
		}); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", ui.WarningText.Render("⚠️  Warning: failed to remove keychain secrets: "+err.Error()))
		}
		
		// Remove from git credential helper (osxkeychain, credential store, etc.)
		if identity.GitHubUser != "" {
			if err := ui.SpinnerWithFunc("Removing git credentials", func() error {
				return keychain.RemoveGitCredentials(identity.GitHubUser)
			}); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", ui.WarningText.Render("⚠️  Warning: failed to remove git credentials: "+err.Error()))
			}
		}
	}

	// Delete SSH key files if requested
	if deleteKeys && identity.SSHKeyPath != "" {
		keyPath := identity.SSHKeyPath
		pubKeyPath := keyPath + ".pub"
		
		if err := os.Remove(keyPath); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%sWarning: failed to delete SSH key: %v%s\n", colorYellow, err, colorReset)
		} else if err == nil {
			fmt.Printf("%s✓ SSH key deleted: %s%s\n", colorGreen, keyPath, colorReset)
		}
		
		if err := os.Remove(pubKeyPath); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%sWarning: failed to delete SSH public key: %v%s\n", colorYellow, err, colorReset)
		} else if err == nil {
			fmt.Printf("%s✓ SSH public key deleted: %s%s\n", colorGreen, pubKeyPath, colorReset)
		}
	}

	// Remove from config
	if err := config.RemoveIdentity(alias); err != nil {
		return fmt.Errorf("failed to remove from config: %w", err)
	}

	fmt.Println()
	fmt.Println(ui.Celebration(fmt.Sprintf("Identity '%s' removed successfully", alias)))
	return nil
}
