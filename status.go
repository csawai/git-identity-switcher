package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/csawai/git-identity-switcher/internal/config"
	"github.com/csawai/git-identity-switcher/internal/ui"
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

	// Get user.name (check local first, then global)
	name, err := getGitConfigLocal("user.name")
	if err != nil {
		name, _ = getGitConfig("user.name")
		if name == "" {
			name = "(not set)"
		}
	}

	// Get user.email (check local first, then global)
	email, err := getGitConfigLocal("user.email")
	if err != nil {
		email, _ = getGitConfig("user.email")
		if email == "" {
			email = "(not set)"
		}
	}

	// Get remote URL
	remote, err := getRemoteURL()
	if err != nil {
		remote = "(not set)"
	}

	// Check if bound to an identity
	boundIdentity := ""

	// First, check for gitx binding marker (most reliable for both SSH and HTTPS)
	marker, err := getGitConfigLocal("gitx.bound")
	if err == nil && marker != "" {
		// Verify the marker matches a stored identity
		cfg, err := config.LoadConfig()
		if err == nil {
			for _, id := range cfg.Identities {
				if id.Alias == marker {
					boundIdentity = marker
					break
				}
			}
		}
	}

	// Fallback: Check if remote uses SSH host alias (for SSH auth)
	// This is the definitive way to detect a bound identity (gitx always sets this for SSH)
	if boundIdentity == "" && strings.Contains(remote, "@github.com-") {
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

	// Build status display
	var statusIcon string
	var statusText string
	var boxStyle lipgloss.Style

	if boundIdentity != "" {
		statusIcon = ui.StatusBound
		statusText = fmt.Sprintf("Bound to: %s", boundIdentity)
		boxStyle = ui.SuccessBox
	} else {
		statusIcon = ui.StatusUnbound
		statusText = "Not bound to any identity"
		boxStyle = ui.WarningBox
	}

	content := fmt.Sprintf(`%s Repository Identity Status

📝 Name:    %s
📧 Email:   %s
🔗 Remote:  %s
%s %s`,
		statusIcon,
		ui.InfoText.Render(name),
		ui.InfoText.Render(email),
		ui.MutedText.Render(remote),
		statusIcon,
		statusText,
	)

	fmt.Println(boxStyle.Render(content))
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

func getGitConfigLocal(key string) (string, error) {
	cmd := exec.Command("git", "config", "--local", "--get", key)
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
