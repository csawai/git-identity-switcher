package config

import (
	"os"
	"testing"
)

func TestConfigDir(t *testing.T) {
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("Failed to get config dir: %v", err)
	}
	if dir == "" {
		t.Error("Config dir should not be empty")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "gitx-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override config dir for test
	originalGetConfigDir := getConfigDirFunc
	defer func() { getConfigDirFunc = originalGetConfigDir }()

	getConfigDirFunc = func() (string, error) {
		return tmpDir, nil
	}

	config := &Config{
		Identities: []Identity{
			{
				Alias:      "test",
				Name:       "Test User",
				Email:      "test@example.com",
				GitHubUser: "testuser",
				AuthMethod: "ssh",
			},
		},
	}

	if err := SaveConfig(config); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(loaded.Identities) != 1 {
		t.Errorf("Expected 1 identity, got %d", len(loaded.Identities))
	}

	if loaded.Identities[0].Alias != "test" {
		t.Errorf("Expected alias 'test', got '%s'", loaded.Identities[0].Alias)
	}
}

func TestAddIdentity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitx-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalGetConfigDir := getConfigDirFunc
	defer func() { getConfigDirFunc = originalGetConfigDir }()

	getConfigDirFunc = func() (string, error) {
		return tmpDir, nil
	}

	identity := Identity{
		Alias:      "work",
		Name:       "Work User",
		Email:      "work@example.com",
		GitHubUser: "workuser",
		AuthMethod: "ssh",
	}

	if err := AddIdentity(identity); err != nil {
		t.Fatalf("Failed to add identity: %v", err)
	}

	found, err := FindIdentityByAlias("work")
	if err != nil {
		t.Fatalf("Failed to find identity: %v", err)
	}

	if found.Email != "work@example.com" {
		t.Errorf("Expected email 'work@example.com', got '%s'", found.Email)
	}
}

func TestRemoveIdentity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitx-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalGetConfigDir := getConfigDirFunc
	defer func() { getConfigDirFunc = originalGetConfigDir }()

	getConfigDirFunc = func() (string, error) {
		return tmpDir, nil
	}

	identity := Identity{
		Alias:      "temp",
		Name:       "Temp User",
		Email:      "temp@example.com",
		GitHubUser: "tempuser",
		AuthMethod: "ssh",
	}

	AddIdentity(identity)

	if err := RemoveIdentity("temp"); err != nil {
		t.Fatalf("Failed to remove identity: %v", err)
	}

	_, err = FindIdentityByAlias("temp")
	if err == nil {
		t.Error("Expected error when finding removed identity")
	}
}

