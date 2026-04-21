// Package profile manages named configuration profiles for vaultpull,
// allowing users to switch between different Vault environments (e.g. dev, staging, prod).
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Profile holds the configuration for a named Vault environment.
type Profile struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Path    string `json:"path"`
	Output  string `json:"output,omitempty"`
}

// Store holds a collection of named profiles.
type Store struct {
	Profiles map[string]Profile `json:"profiles"`
}

// Load reads a profile store from the given file path.
// If the file does not exist, an empty store is returned.
func Load(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Store{Profiles: make(map[string]Profile)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("profile: read %s: %w", path, err)
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("profile: parse %s: %w", path, err)
	}
	if s.Profiles == nil {
		s.Profiles = make(map[string]Profile)
	}
	return &s, nil
}

// Save writes the store to the given file path, creating parent directories as needed.
func Save(path string, s *Store) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("profile: mkdir %s: %w", filepath.Dir(path), err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("profile: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("profile: write %s: %w", path, err)
	}
	return nil
}

// Get retrieves a profile by name. Returns an error if not found.
func (s *Store) Get(name string) (Profile, error) {
	p, ok := s.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile: %q not found", name)
	}
	return p, nil
}

// Set adds or replaces a profile in the store.
func (s *Store) Set(p Profile) error {
	if p.Name == "" {
		return errors.New("profile: name must not be empty")
	}
	if p.Address == "" {
		return errors.New("profile: address must not be empty")
	}
	if p.Path == "" {
		return errors.New("profile: path must not be empty")
	}
	s.Profiles[p.Name] = p
	return nil
}

// Delete removes a profile by name. Returns an error if it does not exist.
func (s *Store) Delete(name string) error {
	if _, ok := s.Profiles[name]; !ok {
		return fmt.Errorf("profile: %q not found", name)
	}
	delete(s.Profiles, name)
	return nil
}
