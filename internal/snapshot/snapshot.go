// Package snapshot provides functionality to save and load local snapshots
// of secret state, enabling rollback and change detection between syncs.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single secret key-value pair captured at a point in time.
type Entry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CapturedAt time.Time `json:"captured_at"`
}

// Snapshot holds a collection of secret entries for a given vault path.
type Snapshot struct {
	VaultPath string    `json:"vault_path"`
	CreatedAt time.Time `json:"created_at"`
	Entries   []Entry   `json:"entries"`
}

// Save writes a snapshot to disk at the given file path.
func Save(path string, s Snapshot) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("snapshot: create dirs: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from disk at the given file path.
func Load(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, nil
		}
		return Snapshot{}, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode: %w", err)
	}
	return s, nil
}

// FromMap builds a Snapshot from a key-value map.
func FromMap(vaultPath string, secrets map[string]string) Snapshot {
	now := time.Now().UTC()
	entries := make([]Entry, 0, len(secrets))
	for k, v := range secrets {
		entries = append(entries, Entry{Key: k, Value: v, CapturedAt: now})
	}
	return Snapshot{VaultPath: vaultPath, CreatedAt: now, Entries: entries}
}

// ToMap converts a Snapshot back into a key-value map.
func ToMap(s Snapshot) map[string]string {
	m := make(map[string]string, len(s.Entries))
	for _, e := range s.Entries {
		m[e.Key] = e.Value
	}
	return m
}
