package keychain

import (
	"fmt"

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

