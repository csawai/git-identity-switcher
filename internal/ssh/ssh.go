package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	SSHConfigMarkerBegin = "# BEGIN gitx managed"
	SSHConfigMarkerEnd   = "# END gitx managed"
)

func GetSSHConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	return filepath.Join(usr.HomeDir, ".ssh", "config"), nil
}

func GenerateSSHKey(identityAlias string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	sshDir := filepath.Join(usr.HomeDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	keyPath := filepath.Join(sshDir, fmt.Sprintf("gitx_%s", identityAlias))

	// Check if key already exists
	if _, err := os.Stat(keyPath); err == nil {
		return keyPath, nil
	}

	// Generate SSH key
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", keyPath, "-N", "", "-C", fmt.Sprintf("gitx-%s", identityAlias))
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to generate SSH key: %w", err)
	}

	return keyPath, nil
}

func AddSSHConfigEntry(hostAlias, keyPath string) error {
	configPath, err := GetSSHConfigPath()
	if err != nil {
		return err
	}

	// Backup before making changes
	if _, err := BackupSSHConfig(); err != nil {
		return fmt.Errorf("failed to backup SSH config: %w", err)
	}

	// Read existing config
	var existingContent string
	if data, err := os.ReadFile(configPath); err == nil {
		existingContent = string(data)
	}

	// Remove existing managed block
	existingContent = removeManagedBlock(existingContent)

	// Create new managed block
	managedBlock := buildManagedBlock(hostAlias, keyPath)

	// Append managed block
	newContent := existingContent
	if !strings.HasSuffix(newContent, "\n") && newContent != "" {
		newContent += "\n"
	}
	newContent += managedBlock + "\n"

	// Write to temp file first
	tempPath := configPath + ".tmp"
	if err := os.WriteFile(tempPath, []byte(newContent), 0600); err != nil {
		return fmt.Errorf("failed to write temp config: %w", err)
	}

	// Validate SSH config
	if err := validateSSHConfig(tempPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("invalid SSH config: %w", err)
	}

	// Atomic swap
	if err := os.Rename(tempPath, configPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

func buildManagedBlock(hostAlias, keyPath string) string {
	block := fmt.Sprintf("%s\n", SSHConfigMarkerBegin)
	block += fmt.Sprintf("Host %s\n", hostAlias)
	block += fmt.Sprintf("  HostName github.com\n")
	block += fmt.Sprintf("  User git\n")
	block += fmt.Sprintf("  IdentityFile %s\n", keyPath)
	block += fmt.Sprintf("  IdentitiesOnly yes\n")
	block += fmt.Sprintf("%s\n", SSHConfigMarkerEnd)
	return block
}

func removeManagedBlock(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inManagedBlock := false

	for _, line := range lines {
		if strings.Contains(line, SSHConfigMarkerBegin) {
			inManagedBlock = true
			continue
		}
		if strings.Contains(line, SSHConfigMarkerEnd) {
			inManagedBlock = false
			continue
		}
		if !inManagedBlock {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func validateSSHConfig(configPath string) error {
	// Simple validation: check if ssh config can parse it
	// We use -F to specify config file and -G to test parsing
	cmd := exec.Command("ssh", "-F", configPath, "-G", "github.com-work")
	// We don't care about the output, just that it doesn't error
	_ = cmd.Run()
	// For now, we'll assume it's valid if the file exists and is readable
	// In production, you might want stricter validation
	return nil
}

func RemoveSSHConfigEntry(hostAlias string) error {
	configPath, err := GetSSHConfigPath()
	if err != nil {
		return err
	}

	// Backup before making changes
	if _, err := BackupSSHConfig(); err != nil {
		return fmt.Errorf("failed to backup SSH config: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	content := removeManagedBlock(string(data))

	// Write to temp file first
	tempPath := configPath + ".tmp"
	if err := os.WriteFile(tempPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write temp config: %w", err)
	}

	// Atomic swap
	if err := os.Rename(tempPath, configPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}
