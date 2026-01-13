package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigDirName  = ".config"
	GitxDirName    = "gitx"
	IdentitiesFile = "identities.json"
)

type Identity struct {
	Alias        string `json:"alias"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	GitHubUser   string `json:"github_user"`
	SSHKeyPath   string `json:"ssh_key_path,omitempty"`
	AuthMethod   string `json:"auth_method"` // "ssh" or "pat"
	SSHHostAlias string `json:"ssh_host_alias,omitempty"`
}

type Config struct {
	Identities []Identity `json:"identities"`
}

var getConfigDirFunc = func() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ConfigDirName, GitxDirName), nil
}

func GetConfigDir() (string, error) {
	return getConfigDirFunc()
}

func GetIdentitiesPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, IdentitiesFile), nil
}

func LoadConfig() (*Config, error) {
	path, err := GetIdentitiesPath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, return empty config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{Identities: []Identity{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	path, err := GetIdentitiesPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func FindIdentityByAlias(alias string) (*Identity, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	for i := range config.Identities {
		if config.Identities[i].Alias == alias {
			return &config.Identities[i], nil
		}
	}

	return nil, fmt.Errorf("identity '%s' not found", alias)
}

func AddIdentity(identity Identity) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// Check if alias already exists
	for _, existing := range config.Identities {
		if existing.Alias == identity.Alias {
			return fmt.Errorf("identity with alias '%s' already exists", identity.Alias)
		}
	}

	config.Identities = append(config.Identities, identity)
	return SaveConfig(config)
}

func RemoveIdentity(alias string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	found := false
	identities := []Identity{}
	for _, id := range config.Identities {
		if id.Alias == alias {
			found = true
		} else {
			identities = append(identities, id)
		}
	}

	if !found {
		return fmt.Errorf("identity '%s' not found", alias)
	}

	config.Identities = identities
	return SaveConfig(config)
}

