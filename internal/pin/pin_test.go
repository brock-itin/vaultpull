package pin_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vaultpull/internal/pin"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")

	pf := &pin.PinFile{Pins: make(map[string]pin.Entry)}
	pf.Set("secret/app/db", 3, "alice")
	pf.Set("secret/app/api", 1, "bob")

	if err := pin.Save(path, pf); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := pin.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Pins) != 2 {
		t.Fatalf("expected 2 pins, got %d", len(loaded.Pins))
	}
	e, ok := loaded.Get("secret/app/db")
	if !ok {
		t.Fatal("expected pin for secret/app/db")
	}
	if e.Version != 3 {
		t.Errorf("expected version 3, got %d", e.Version)
	}
	if e.PinnedBy != "alice" {
		t.Errorf("expected pinnedBy alice, got %s", e.PinnedBy)
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	pf, err := pin.Load("/nonexistent/pins.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pf.Pins) != 0 {
		t.Errorf("expected empty pins, got %d", len(pf.Pins))
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "pins.json")
	pf := &pin.PinFile{Pins: make(map[string]pin.Entry)}
	pf.Set("secret/x", 1, "")
	if err := pin.Save(path, pf); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	pf := &pin.PinFile{Pins: make(map[string]pin.Entry)}
	pf.Set("secret/a", 2, "")
	pf.Remove("secret/a")
	if _, ok := pf.Get("secret/a"); ok {
		t.Error("expected pin to be removed")
	}
}

func TestCheckDrift_DetectsMismatch(t *testing.T) {
	pf := &pin.PinFile{Pins: make(map[string]pin.Entry)}
	pf.Set("secret/a", 2, "")
	pf.Set("secret/b", 5, "")

	current := map[string]int{
		"secret/a": 2, // same — no drift
		"secret/b": 6, // different — drift
	}
	drifted := pf.CheckDrift(current)
	if len(drifted) != 1 || drifted[0] != "secret/b" {
		t.Errorf("expected [secret/b] drifted, got %v", drifted)
	}
}

func TestCheckDrift_MissingPathIsDrift(t *testing.T) {
	pf := &pin.PinFile{Pins: make(map[string]pin.Entry)}
	pf.Set("secret/gone", 1, "")
	drifted := pf.CheckDrift(map[string]int{})
	if len(drifted) != 1 {
		t.Errorf("expected 1 drifted path, got %v", drifted)
	}
}

func TestSet_UpdatesPinnedAt(t *testing.T) {
	pf := &pin.PinFile{Pins: make(map[string]pin.Entry)}
	before := time.Now().UTC()
	pf.Set("secret/ts", 1, "")
	e, _ := pf.Get("secret/ts")
	if e.PinnedAt.Before(before) {
		t.Error("PinnedAt should be >= before")
	}
}
