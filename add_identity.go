package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/csawai/gitx/internal/config"
	"github.com/csawai/gitx/internal/keychain"
	"github.com/csawai/gitx/internal/ssh"
	"github.com/spf13/cobra"
)

var dryRun bool

func init() {
	addIdentityCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")
}

var addIdentityCmd = &cobra.Command{
	Use:   "add identity",
	Short: "Add a new identity",
	Long:  "Add a new GitHub identity with name, email, and GitHub username.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := addIdentity(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func addIdentity() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Identity alias (e.g., 'work', 'personal'): ")
	alias, _ := reader.ReadString('\n')
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return fmt.Errorf("alias cannot be empty")
	}

	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	fmt.Print("GitHub username: ")
	githubUser, _ := reader.ReadString('\n')
	githubUser = strings.TrimSpace(githubUser)
	if githubUser == "" {
		return fmt.Errorf("GitHub username cannot be empty")
	}

	// Ask for auth method
	fmt.Print("Auth method (ssh/pat) [ssh]: ")
	authMethod, _ := reader.ReadString('\n')
	authMethod = strings.TrimSpace(authMethod)
	if authMethod == "" {
		authMethod = "ssh"
	}

	identity := config.Identity{
		Alias:      alias,
		Name:       name,
		Email:      email,
		GitHubUser: githubUser,
		AuthMethod: authMethod,
	}

	if dryRun {
		fmt.Println("\n[DRY RUN] Would add identity:")
		fmt.Printf("  Alias: %s\n", alias)
		fmt.Printf("  Name: %s\n", name)
		fmt.Printf("  Email: %s\n", email)
		fmt.Printf("  GitHub: %s\n", githubUser)
		fmt.Printf("  Auth: %s\n", authMethod)
		return nil
	}

	// Handle SSH key generation
	if authMethod == "ssh" {
		fmt.Print("Generate SSH key? (y/n) [y]: ")
		generate, _ := reader.ReadString('\n')
		generate = strings.TrimSpace(strings.ToLower(generate))
		if generate == "" || generate == "y" {
			// Backup SSH config first
			backupPath, err := ssh.BackupSSHConfig()
			if err != nil {
				return fmt.Errorf("failed to backup SSH config: %w", err)
			}
			if backupPath != "" {
				fmt.Printf("✓ SSH config backed up to: %s\n", backupPath)
			}

			keyPath, err := ssh.GenerateSSHKey(alias)
			if err != nil {
				return fmt.Errorf("failed to generate SSH key: %w", err)
			}
			identity.SSHKeyPath = keyPath
			identity.SSHHostAlias = fmt.Sprintf("github.com-%s", alias)

			// Add SSH config entry
			if err := ssh.AddSSHConfigEntry(identity.SSHHostAlias, keyPath); err != nil {
				return fmt.Errorf("failed to add SSH config: %w", err)
			}
			fmt.Printf("✓ SSH key generated: %s\n", keyPath)
			fmt.Printf("✓ SSH config updated\n")
		}
	} else if authMethod == "pat" {
		fmt.Print("Personal Access Token: ")
		token, _ := reader.ReadString('\n')
		token = strings.TrimSpace(token)
		if token == "" {
			return fmt.Errorf("PAT cannot be empty")
		}

		// Store PAT in keychain
		if err := keychain.StoreSecret(alias, "pat", token); err != nil {
			return fmt.Errorf("failed to store PAT: %w", err)
		}
		fmt.Println("✓ PAT stored securely in keychain")
	}

	if err := config.AddIdentity(identity); err != nil {
		return err
	}

	fmt.Printf("✓ Identity '%s' added successfully\n", alias)
	return nil
}
