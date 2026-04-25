// Package pin provides version pinning for Vault secret paths,
// allowing callers to lock a path to a specific version and detect drift.
package pin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry records a pinned version for a single Vault path.
type Entry struct {
	Path      string    `json:"path"`
	Version   int       `json:"version"`
	PinnedAt  time.Time `json:"pinned_at"`
	PinnedBy  string    `json:"pinned_by,omitempty"`
}

// PinFile is the top-level structure persisted to disk.
type PinFile struct {
	Pins map[string]Entry `json:"pins"`
}

// Load reads a pin file from disk. If the file does not exist an empty
// PinFile is returned without error.
func Load(path string) (*PinFile, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &PinFile{Pins: make(map[string]Entry)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("pin: read %s: %w", path, err)
	}
	var pf PinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("pin: parse %s: %w", path, err)
	}
	if pf.Pins == nil {
		pf.Pins = make(map[string]Entry)
	}
	return &pf, nil
}

// Save writes the PinFile to disk, creating parent directories as needed.
func Save(path string, pf *PinFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("pin: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("pin: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("pin: write %s: %w", path, err)
	}
	return nil
}

// Set adds or updates a pin entry for the given path.
func (pf *PinFile) Set(vaultPath string, version int, pinnedBy string) {
	pf.Pins[vaultPath] = Entry{
		Path:     vaultPath,
		Version:  version,
		PinnedAt: time.Now().UTC(),
		PinnedBy: pinnedBy,
	}
}

// Get returns the pinned entry for a path and whether it exists.
func (pf *PinFile) Get(vaultPath string) (Entry, bool) {
	e, ok := pf.Pins[vaultPath]
	return e, ok
}

// Remove deletes the pin for a path if present.
func (pf *PinFile) Remove(vaultPath string) {
	delete(pf.Pins, vaultPath)
}

// CheckDrift returns paths whose current version differs from the pinned version.
// current is a map of vaultPath -> observed version.
func (pf *PinFile) CheckDrift(current map[string]int) []string {
	var drifted []string
	for path, entry := range pf.Pins {
		v, ok := current[path]
		if !ok || v != entry.Version {
			drifted = append(drifted, path)
		}
	}
	return drifted
}
