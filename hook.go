package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var installHookCmd = &cobra.Command{
	Use:   "install-hook",
	Short: "Install pre-push hook",
	Long:  "Install a pre-push hook that blocks pushes when repository is not bound to an identity.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := installHook(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var uninstallHookCmd = &cobra.Command{
	Use:   "uninstall-hook",
	Short: "Uninstall pre-push hook",
	Long:  "Remove the gitx pre-push hook.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := uninstallHook(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installHookCmd)
	rootCmd.AddCommand(uninstallHookCmd)
}

func installHook() error {
	// Check if we're in a git repo
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("not a git repository")
	}

	gitDir := strings.TrimSpace(string(output))
	hooksDir := filepath.Join(gitDir, "hooks")
	hookPath := filepath.Join(hooksDir, "pre-push")

	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		// Check if it's our hook
		data, _ := os.ReadFile(hookPath)
		if strings.Contains(string(data), "gitx") {
			fmt.Println("gitx pre-push hook already installed")
			return nil
		}
		// Existing hook - we should merge or warn
		fmt.Println("Warning: pre-push hook already exists. gitx hook not installed.")
		fmt.Println("You can manually merge the gitx check into your existing hook.")
		return nil
	}

	// Create hooks directory if it doesn't exist
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Write hook
	hookContent := `#!/bin/sh
# gitx pre-push hook
# Blocks push if repository is not bound to an identity

# Check if remote uses gitx host alias
remote=$(git remote get-url origin 2>/dev/null)
if [ -z "$remote" ]; then
  echo "Error: No remote configured"
  exit 1
fi

# Check if bound (contains github.com-)
if echo "$remote" | grep -q "github.com-"; then
  exit 0
fi

# Check if user.name and user.email are set locally
name=$(git config --local --get user.name 2>/dev/null)
email=$(git config --local --get user.email 2>/dev/null)

if [ -n "$name" ] && [ -n "$email" ]; then
  exit 0
fi

echo "Error: Repository is not bound to an identity."
echo "Run 'gitx bind <identity>' to bind this repository."
exit 1
`

	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		return fmt.Errorf("failed to write hook: %w", err)
	}

	fmt.Println("✓ Pre-push hook installed")
	return nil
}

func uninstallHook() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("not a git repository")
	}

	gitDir := strings.TrimSpace(string(output))
	hookPath := filepath.Join(gitDir, "hooks", "pre-push")

	// Check if hook exists and is ours
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		fmt.Println("No gitx pre-push hook found")
		return nil
	}

	data, err := os.ReadFile(hookPath)
	if err != nil {
		return err
	}

	if !strings.Contains(string(data), "gitx") {
		fmt.Println("Pre-push hook exists but is not a gitx hook")
		return nil
	}

	if err := os.Remove(hookPath); err != nil {
		return fmt.Errorf("failed to remove hook: %w", err)
	}

	fmt.Println("✓ Pre-push hook uninstalled")
	return nil
}

