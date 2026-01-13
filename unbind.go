package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/csawai/git-identity-switcher/internal/ui"
	"github.com/spf13/cobra"
)

var unbindCmd = &cobra.Command{
	Use:   "unbind",
	Short: "Unbind repository from identity",
	Long:  "Revert repository-local git config changes made by gitx.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := unbind(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func unbind() error {
	// Check if we're in a git repo
	gitCmd := exec.Command("git", "rev-parse", "--git-dir")
	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("not a git repository")
	}

	// Unset user.name if it was set by gitx
	// For now, we'll just unset it (in future, we could track original values)
	unsetName := exec.Command("git", "config", "--local", "--unset", "user.name")
	if err := unsetName.Run(); err != nil {
		// Key might not exist, that's okay
	}

	unsetEmail := exec.Command("git", "config", "--local", "--unset", "user.email")
	if err := unsetEmail.Run(); err != nil {
		// Key might not exist, that's okay
	}

	// Verify they're actually unset (double-check)
	verifyName := exec.Command("git", "config", "--local", "--get", "user.name")
	if verifyName.Run() == nil {
		// Still exists, try unset-all to remove all occurrences
		exec.Command("git", "config", "--local", "--unset-all", "user.name").Run()
	}

	verifyEmail := exec.Command("git", "config", "--local", "--get", "user.email")
	if verifyEmail.Run() == nil {
		// Still exists, try unset-all to remove all occurrences
		exec.Command("git", "config", "--local", "--unset-all", "user.email").Run()
	}

	// Remove gitx binding marker
	unsetMarker := exec.Command("git", "config", "--local", "--unset", "gitx.bound")
	_ = unsetMarker.Run() // Ignore error if not set

	// Try to revert remote URL to standard github.com format
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err == nil {
		currentURL := strings.TrimSpace(string(output))
		// If using a host alias, convert back to standard format
		if strings.Contains(currentURL, "@github.com-") {
			// Extract org/repo
			// Format: git@github.com-IDENTITY:org/repo.git
			parts := strings.Split(currentURL, ":")
			if len(parts) == 2 {
				newURL := fmt.Sprintf("git@github.com:%s", parts[1])
				setRemoteCmd := exec.Command("git", "remote", "set-url", "origin", newURL)
				if err := setRemoteCmd.Run(); err != nil {
					// Log error but don't fail - remote URL revert is best effort
					fmt.Fprintf(os.Stderr, "Warning: could not revert remote URL: %v\n", err)
				}
			}
		}
	}

	fmt.Println(ui.SuccessBox.Render("✅ Repository unbound successfully"))
	return nil
}
