package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func BackupSSHConfig() (string, error) {
	configPath, err := GetSSHConfigPath()
	if err != nil {
		return "", err
	}

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "", nil // No backup needed if file doesn't exist
	}

	// Create backup with timestamp
	backupPath := fmt.Sprintf("%s.gitx.backup.%s", configPath, time.Now().Format("20060102-150405"))
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read config for backup: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return "", fmt.Errorf("failed to write backup: %w", err)
	}

	return backupPath, nil
}

func CleanupOldBackups(keepCount int) error {
	configPath, err := GetSSHConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	pattern := filepath.Join(configDir, "config.gitx.backup.*")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= keepCount {
		return nil
	}

	// Sort by modification time (oldest first)
	// Simple approach: delete all but the most recent keepCount
	for i := 0; i < len(matches)-keepCount; i++ {
		os.Remove(matches[i])
	}

	return nil
}

