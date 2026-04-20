// Package lineage tracks the origin of secrets pulled from Vault,
// recording which path and version each key was sourced from.
package lineage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry records where a single secret key was sourced from.
type Entry struct {
	Key       string    `json:"key"`
	VaultPath string    `json:"vault_path"`
	Version   int       `json:"version"`
	FetchedAt time.Time `json:"fetched_at"`
}

// Record is a collection of lineage entries keyed by secret key name.
type Record map[string]Entry

// Build constructs a Record from a map of secret key→value pairs,
// annotating each with the given Vault path, version, and current time.
func Build(secrets map[string]string, vaultPath string, version int) Record {
	now := time.Now().UTC()
	r := make(Record, len(secrets))
	for k := range secrets {
		r[k] = Entry{
			Key:       k,
			VaultPath: vaultPath,
			Version:   version,
			FetchedAt: now,
		}
	}
	return r
}

// Save writes the Record as JSON to the given file path,
// creating any necessary parent directories.
func Save(path string, r Record) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// Load reads a Record from a JSON file. If the file does not exist,
// an empty Record is returned without error.
func Load(path string) (Record, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Record{}, nil
	}
	if err != nil {
		return nil, err
	}
	var r Record
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return r, nil
}
