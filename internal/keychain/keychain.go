package keychain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/99designs/keyring"
)

const ServiceName = "gitx"

func getKeyring() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: ServiceName,
	})
}

func StoreSecret(identityAlias, key, value string) error {
	ring, err := getKeyring()
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}
	secretKey := fmt.Sprintf("%s:%s", identityAlias, key)
	return ring.Set(keyring.Item{
		Key:  secretKey,
		Data: []byte(value),
	})
}

func GetSecret(identityAlias, key string) (string, error) {
	ring, err := getKeyring()
	if err != nil {
		return "", fmt.Errorf("failed to open keyring: %w", err)
	}
	secretKey := fmt.Sprintf("%s:%s", identityAlias, key)
	item, err := ring.Get(secretKey)
	if err != nil {
		return "", err
	}
	return string(item.Data), nil
}

func DeleteSecret(identityAlias, key string) error {
	ring, err := getKeyring()
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}
	secretKey := fmt.Sprintf("%s:%s", identityAlias, key)
	return ring.Remove(secretKey)
}

func DeleteAllSecrets(identityAlias string) error {
	// Note: keyring doesn't support listing easily, so we delete known keys
	keys := []string{"pat", "ssh_passphrase"}
	for _, key := range keys {
		DeleteSecret(identityAlias, key) // Ignore errors
	}
	return nil
}

// RemoveGitCredentials removes credentials from git's credential helper
// This handles both osxkeychain (macOS) and credential store (Linux)
func RemoveGitCredentials(githubUser string) error {
	// Try osxkeychain first (macOS)
	cmd := exec.Command("git", "credential-osxkeychain", "erase")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("protocol=https\nhost=github.com\nusername=%s\n\n", githubUser))
	if err := cmd.Run(); err == nil {
		return nil // Success
	}

	// Try credential store (Linux/Windows)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	credStorePath := filepath.Join(homeDir, ".git-credentials")
	if _, err := os.Stat(credStorePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to remove
	}

	data, err := os.ReadFile(credStorePath)
	if err != nil {
		return fmt.Errorf("failed to read credential store: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var filteredLines []string
	pattern := fmt.Sprintf("https://%s@github.com", githubUser)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, pattern) {
			filteredLines = append(filteredLines, line)
		}
	}

	// Also remove generic github.com entries if they match
	// (Some users might have stored without username)
	newContent := strings.Join(filteredLines, "\n")
	if newContent != "" && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}

	// Write back
	if err := os.WriteFile(credStorePath, []byte(newContent), 0600); err != nil {
		return fmt.Errorf("failed to write credential store: %w", err)
	}

	return nil
}

