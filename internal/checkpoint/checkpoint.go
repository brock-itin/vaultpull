// Package checkpoint tracks the last successful sync time and metadata
// for each configured vault path, enabling incremental and resumable syncs.
package checkpoint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry records the result of a single sync operation for one path.
type Entry struct {
	Path      string    `json:"path"`
	SyncedAt  time.Time `json:"synced_at"`
	KeyCount  int       `json:"key_count"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}

// Checkpoint holds sync state for all tracked paths.
type Checkpoint struct {
	Entries map[string]Entry `json:"entries"`
}

// New returns an empty Checkpoint.
func New() *Checkpoint {
	return &Checkpoint{Entries: make(map[string]Entry)}
}

// Record stores or updates the entry for the given vault path.
func (c *Checkpoint) Record(e Entry) {
	c.Entries[e.Path] = e
}

// Get returns the entry for a path and whether it exists.
func (c *Checkpoint) Get(path string) (Entry, bool) {
	e, ok := c.Entries[path]
	return e, ok
}

// Save writes the checkpoint to disk as JSON.
func Save(path string, c *Checkpoint) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}

// Load reads a checkpoint from disk. If the file does not exist, an empty
// Checkpoint is returned without error.
func Load(path string) (*Checkpoint, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return New(), nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c := New()
	if err := json.NewDecoder(f).Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}
