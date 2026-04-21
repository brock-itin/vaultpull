package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/checkpoint"
)

func TestRecord_AndGet(t *testing.T) {
	c := checkpoint.New()
	e := checkpoint.Entry{
		Path:     "secret/app",
		SyncedAt: time.Now().UTC().Truncate(time.Second),
		KeyCount: 5,
		Success:  true,
	}
	c.Record(e)
	got, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.KeyCount != 5 {
		t.Errorf("key count: got %d, want 5", got.KeyCount)
	}
}

func TestGet_Missing(t *testing.T) {
	c := checkpoint.New()
	_, ok := c.Get("secret/missing")
	if ok {
		t.Fatal("expected no entry for unknown path")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "checkpoint.json")

	c := checkpoint.New()
	c.Record(checkpoint.Entry{
		Path:     "secret/db",
		SyncedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		KeyCount: 3,
		Success:  true,
	})

	if err := checkpoint.Save(path, c); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := checkpoint.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	e, ok := loaded.Get("secret/db")
	if !ok {
		t.Fatal("entry not found after load")
	}
	if e.KeyCount != 3 {
		t.Errorf("key count: got %d, want 3", e.KeyCount)
	}
	if !e.Success {
		t.Error("expected success=true")
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	c, err := checkpoint.Load(filepath.Join(dir, "no-such-file.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Entries) != 0 {
		t.Errorf("expected empty checkpoint, got %d entries", len(c.Entries))
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "cp.json")
	if err := checkpoint.Save(path, checkpoint.New()); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestSave_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cp.json")
	if err := checkpoint.Save(path, checkpoint.New()); err != nil {
		t.Fatalf("Save: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("permissions: got %o, want 0600", perm)
	}
}
