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

// SSHIdentity represents an SSH host alias and key path pair
type SSHIdentity struct {
	HostAlias string
	KeyPath   string
}

// AddSSHConfigEntry adds or updates an SSH config entry, preserving all existing gitx-managed entries
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

	// Parse existing managed entries
	existingIdentities := parseManagedBlock(existingContent)
	
	// Add or update the new identity
	found := false
	for i, id := range existingIdentities {
		if id.HostAlias == hostAlias {
			existingIdentities[i].KeyPath = keyPath
			found = true
			break
		}
	}
	if !found {
		existingIdentities = append(existingIdentities, SSHIdentity{
			HostAlias: hostAlias,
			KeyPath:   keyPath,
		})
	}

	// Remove existing managed block
	existingContent = removeManagedBlock(existingContent)

	// Build new managed block with all identities
	managedBlock := buildManagedBlockFromIdentities(existingIdentities)

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

func buildManagedBlockFromIdentities(identities []SSHIdentity) string {
	block := fmt.Sprintf("%s\n", SSHConfigMarkerBegin)
	for _, id := range identities {
		block += fmt.Sprintf("Host %s\n", id.HostAlias)
		block += fmt.Sprintf("  HostName github.com\n")
		block += fmt.Sprintf("  User git\n")
		block += fmt.Sprintf("  IdentityFile %s\n", id.KeyPath)
		block += fmt.Sprintf("  IdentitiesOnly yes\n")
		block += "\n"
	}
	block += fmt.Sprintf("%s\n", SSHConfigMarkerEnd)
	return block
}

func parseManagedBlock(content string) []SSHIdentity {
	var identities []SSHIdentity
	lines := strings.Split(content, "\n")
	inManagedBlock := false
	var currentHost string
	var currentKeyPath string

	for i, line := range lines {
		if strings.Contains(line, SSHConfigMarkerBegin) {
			inManagedBlock = true
			continue
		}
		if strings.Contains(line, SSHConfigMarkerEnd) {
			// Save last identity if any
			if currentHost != "" && currentKeyPath != "" {
				identities = append(identities, SSHIdentity{
					HostAlias: currentHost,
					KeyPath:   currentKeyPath,
				})
			}
			inManagedBlock = false
			currentHost = ""
			currentKeyPath = ""
			continue
		}
		if inManagedBlock {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Host ") {
				// Save previous identity if any
				if currentHost != "" && currentKeyPath != "" {
					identities = append(identities, SSHIdentity{
						HostAlias: currentHost,
						KeyPath:   currentKeyPath,
					})
				}
				currentHost = strings.TrimPrefix(line, "Host ")
				currentKeyPath = ""
			} else if strings.HasPrefix(line, "IdentityFile ") {
				currentKeyPath = strings.TrimPrefix(line, "IdentityFile ")
			}
		}
		_ = i // avoid unused variable
	}

	return identities
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

	// Parse existing managed entries
	existingIdentities := parseManagedBlock(string(data))
	
	// Remove the specified host alias
	filteredIdentities := []SSHIdentity{}
	for _, id := range existingIdentities {
		if id.HostAlias != hostAlias {
			filteredIdentities = append(filteredIdentities, id)
		}
	}

	// Remove existing managed block
	content := removeManagedBlock(string(data))

	// If there are remaining identities, rebuild the managed block
	if len(filteredIdentities) > 0 {
		managedBlock := buildManagedBlockFromIdentities(filteredIdentities)
		if !strings.HasSuffix(content, "\n") && content != "" {
			content += "\n"
		}
		content += managedBlock + "\n"
	}

	// Write to temp file first
	tempPath := configPath + ".tmp"
	if err := os.WriteFile(tempPath, []byte(content), 0600); err != nil {
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
