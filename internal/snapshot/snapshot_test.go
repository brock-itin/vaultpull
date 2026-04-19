package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/snapshot"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.Snapshot{
		VaultPath: "secret/myapp",
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		Entries: []snapshot.Entry{
			{Key: "DB_PASS", Value: "s3cr3t", CapturedAt: time.Now().UTC().Truncate(time.Second)},
		},
	}

	if err := snapshot.Save(path, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.VaultPath != orig.VaultPath {
		t.Errorf("VaultPath: got %q want %q", loaded.VaultPath, orig.VaultPath)
	}
	if len(loaded.Entries) != 1 || loaded.Entries[0].Key != "DB_PASS" {
		t.Errorf("unexpected entries: %+v", loaded.Entries)
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	s, err := snapshot.Load("/nonexistent/path/snap.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if s.VaultPath != "" || len(s.Entries) != 0 {
		t.Errorf("expected empty snapshot, got: %+v", s)
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dir", "snap.json")
	s := snapshot.FromMap("secret/app", map[string]string{"KEY": "val"})
	if err := snapshot.Save(path, s); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestFromMap_ToMap_RoundTrip(t *testing.T) {
	input := map[string]string{"A": "1", "B": "2"}
	s := snapshot.FromMap("secret/test", input)
	if s.VaultPath != "secret/test" {
		t.Errorf("unexpected VaultPath: %s", s.VaultPath)
	}
	out := snapshot.ToMap(s)
	for k, v := range input {
		if out[k] != v {
			t.Errorf("key %q: got %q want %q", k, out[k], v)
		}
	}
}

func TestSave_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	s := snapshot.FromMap("secret/app", map[string]string{"X": "y"})
	if err := snapshot.Save(path, s); err != nil {
		t.Fatalf("Save: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
