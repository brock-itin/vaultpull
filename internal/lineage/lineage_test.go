package lineage_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/lineage"
)

func TestBuild_PopulatesEntries(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "abc"}
	before := time.Now().UTC()
	r := lineage.Build(secrets, "secret/myapp", 3)
	after := time.Now().UTC()

	if len(r) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(r))
	}
	for _, k := range []string{"DB_PASS", "API_KEY"} {
		e, ok := r[k]
		if !ok {
			t.Fatalf("missing entry for %s", k)
		}
		if e.VaultPath != "secret/myapp" {
			t.Errorf("expected vault_path secret/myapp, got %s", e.VaultPath)
		}
		if e.Version != 3 {
			t.Errorf("expected version 3, got %d", e.Version)
		}
		if e.FetchedAt.Before(before) || e.FetchedAt.After(after) {
			t.Errorf("FetchedAt %v out of expected range", e.FetchedAt)
		}
	}
}

func TestBuild_EmptySecrets(t *testing.T) {
	r := lineage.Build(map[string]string{}, "secret/empty", 1)
	if len(r) != 0 {
		t.Errorf("expected empty record, got %d entries", len(r))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lineage.json")

	secrets := map[string]string{"TOKEN": "x", "SECRET": "y"}
	orig := lineage.Build(secrets, "secret/svc", 2)

	if err := lineage.Save(path, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := lineage.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != len(orig) {
		t.Fatalf("expected %d entries, got %d", len(orig), len(loaded))
	}
	for k, e := range orig {
		l, ok := loaded[k]
		if !ok {
			t.Fatalf("missing key %s after load", k)
		}
		if l.VaultPath != e.VaultPath || l.Version != e.Version {
			t.Errorf("mismatch for key %s", k)
		}
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	r, err := lineage.Load("/nonexistent/path/lineage.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(r) != 0 {
		t.Errorf("expected empty record, got %d entries", len(r))
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "deep", "lineage.json")
	r := lineage.Build(map[string]string{"K": "v"}, "secret/p", 1)
	if err := lineage.Save(path, r); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
