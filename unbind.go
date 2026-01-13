package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	exec.Command("git", "config", "--local", "--unset", "user.name").Run()
	exec.Command("git", "config", "--local", "--unset", "user.email").Run()

	// Try to revert remote URL to standard github.com format
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err == nil {
		currentURL := strings.TrimSpace(string(output))
		// If using a host alias, convert back to standard format
		if strings.Contains(currentURL, "@github.com-") {
			// Extract org/repo
			parts := strings.Split(currentURL, ":")
			if len(parts) == 2 {
				newURL := fmt.Sprintf("git@github.com:%s", parts[1])
				exec.Command("git", "remote", "set-url", "origin", newURL).Run()
			}
		}
	}

	fmt.Println("✓ Repository unbound")
	return nil
}

