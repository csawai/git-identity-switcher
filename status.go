package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/csawai/git-identity-switcher/internal/config"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current repository identity status",
	Long:  "Displays the current git user.name, user.email, and remote URL for the repository.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := showStatus(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func showStatus() error {
	// Check if we're in a git repository
	if !isGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	// Get user.name
	name, err := getGitConfig("user.name")
	if err != nil {
		name = "(not set)"
	}

	// Get user.email
	email, err := getGitConfig("user.email")
	if err != nil {
		email = "(not set)"
	}

	// Get remote URL
	remote, err := getRemoteURL()
	if err != nil {
		remote = "(not set)"
	}

	// Check if bound to an identity
	boundIdentity := ""
	
	// Check if remote uses SSH host alias (for SSH auth)
	if strings.Contains(remote, "@github.com-") {
		parts := strings.Split(remote, "@")
		if len(parts) > 1 {
			hostParts := strings.Split(parts[1], ":")
			if len(hostParts) > 0 {
				hostAlias := hostParts[0]
				if strings.HasPrefix(hostAlias, "github.com-") {
					boundIdentity = strings.TrimPrefix(hostAlias, "github.com-")
				}
			}
		}
	}
	
	// If not detected via SSH alias, check if email matches a stored identity
	// (email is the key identifier, name might vary)
	if boundIdentity == "" && email != "(not set)" {
		cfg, err := config.LoadConfig()
		if err == nil {
			for _, id := range cfg.Identities {
				if id.Email == email {
					boundIdentity = id.Alias
					break
				}
			}
		}
	}

	fmt.Println("Repository Identity Status:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Name:  %s\n", name)
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Remote: %s\n", remote)
	if boundIdentity != "" {
		fmt.Printf("Bound to: %s\n", boundIdentity)
	} else {
		fmt.Println("Bound to: (not bound)")
	}

	return nil
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func getGitConfig(key string) (string, error) {
	cmd := exec.Command("git", "config", "--get", key)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getRemoteURL() (string, error) {
	// Try origin first
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	// Fallback to any remote
	cmd = exec.Command("git", "remote", "-v")
	output, err = cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			return parts[1], nil
		}
	}

	return "", fmt.Errorf("no remote found")
}
