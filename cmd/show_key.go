package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/csawai/git-identity-switcher/internal/config"
	"github.com/spf13/cobra"
)

var showKeyCmd = &cobra.Command{
	Use:   "show-key [alias]",
	Short: "Show SSH public key for an identity",
	Long:  "Display the SSH public key that needs to be added to GitHub.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := showKey(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var copyKeyCmd = &cobra.Command{
	Use:   "copy-key [alias]",
	Short: "Copy SSH public key to clipboard",
	Long:  "Copy the SSH public key to clipboard for easy pasting into GitHub.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := copyKey(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(showKeyCmd)
	rootCmd.AddCommand(copyKeyCmd)
}

func showKey(alias string) error {
	identity, err := config.FindIdentityByAlias(alias)
	if err != nil {
		return err
	}

	if identity.SSHKeyPath == "" {
		return fmt.Errorf("identity '%s' does not have an SSH key", alias)
	}

	pubKeyPath := identity.SSHKeyPath + ".pub"
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("SSH public key not found: %s", pubKeyPath)
	}

	pubKey, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("SSH Public Key for '%s':\n", alias)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Print(string(pubKey))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("Add this key to your GitHub account:")
	fmt.Println("  https://github.com/settings/ssh/new")
	fmt.Println()
	fmt.Println("Or copy it with: gitx copy-key", alias)

	return nil
}

func copyKey(alias string) error {
	identity, err := config.FindIdentityByAlias(alias)
	if err != nil {
		return err
	}

	if identity.SSHKeyPath == "" {
		return fmt.Errorf("identity '%s' does not have an SSH key", alias)
	}

	pubKeyPath := identity.SSHKeyPath + ".pub"
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("SSH public key not found: %s", pubKeyPath)
	}

	pubKey, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	// Detect OS and use appropriate clipboard command
	var copyCmd *exec.Cmd
	if _, err := exec.LookPath("pbcopy"); err == nil {
		// macOS
		copyCmd = exec.Command("pbcopy")
	} else if _, err := exec.LookPath("xclip"); err == nil {
		// Linux with xclip
		copyCmd = exec.Command("xclip", "-selection", "clipboard")
	} else if _, err := exec.LookPath("xsel"); err == nil {
		// Linux with xsel
		copyCmd = exec.Command("xsel", "--clipboard", "--input")
	} else {
		return fmt.Errorf("no clipboard utility found (pbcopy, xclip, or xsel)")
	}

	copyCmd.Stdin = os.Stdin
	stdin, err := copyCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := copyCmd.Start(); err != nil {
		return fmt.Errorf("failed to start copy command: %w", err)
	}

	if _, err := stdin.Write(pubKey); err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}
	stdin.Close()

	if err := copyCmd.Wait(); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	fmt.Printf("✓ SSH public key for '%s' copied to clipboard\n", alias)
	fmt.Println("Paste it at: https://github.com/settings/ssh/new")

	return nil
}

