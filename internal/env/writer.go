// Package env provides utilities for writing secrets to .env files.
package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteOptions configures how secrets are written to the .env file.
type WriteOptions struct {
	// Overwrite controls whether existing keys are overwritten.
	Overwrite bool
	// BackupExisting creates a .env.bak before writing.
	BackupExisting bool
}

// Write writes the provided secrets map to the given .env file path.
// Existing entries not present in secrets are preserved.
func Write(path string, secrets map[string]string, opts WriteOptions) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("env: create directories: %w", err)
	}

	existing, err := readExisting(path)
	if err != nil {
		return err
	}

	if opts.BackupExisting && len(existing) > 0 {
		if err := backupFile(path); err != nil {
			return err
		}
	}

	merged := mergeSecrets(existing, secrets, opts.Overwrite)

	return writeFile(path, merged)
}

func readExisting(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("env: read existing file: %w", err)
	}

	result := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result, nil
}

func mergeSecrets(existing, incoming map[string]string, overwrite bool) map[string]string {
	merged := make(map[string]string, len(existing))
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range incoming {
		if _, exists := merged[k]; !exists || overwrite {
			merged[k] = v
		}
	}
	return merged
}

func writeFile(path string, secrets map[string]string) error {
	var sb strings.Builder
	for k, v := range secrets {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
	if err := os.WriteFile(path, []byte(sb.String()), 0o600); err != nil {
		return fmt.Errorf("env: write file: %w", err)
	}
	return nil
}

func backupFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("env: backup read: %w", err)
	}
	if err := os.WriteFile(path+".bak", data, 0o600); err != nil {
		return fmt.Errorf("env: backup write: %w", err)
	}
	return nil
}
