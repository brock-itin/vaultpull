package profile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vaultpull/vaultpull/internal/profile"
)

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.json")

	s := &profile.Store{Profiles: map[string]profile.Profile{
		"dev": {Name: "dev", Address: "http://localhost:8200", Path: "secret/dev", Output: ".env.dev"},
	}}

	if err := profile.Save(path, s); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := profile.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	p, err := loaded.Get("dev")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if p.Address != "http://localhost:8200" {
		t.Errorf("address = %q, want %q", p.Address, "http://localhost:8200")
	}
	if p.Output != ".env.dev" {
		t.Errorf("output = %q, want %q", p.Output, ".env.dev")
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	s, err := profile.Load("/nonexistent/path/profiles.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(s.Profiles) != 0 {
		t.Errorf("expected empty store, got %d profiles", len(s.Profiles))
	}
}

func TestSet_ValidProfile(t *testing.T) {
	s := &profile.Store{Profiles: make(map[string]profile.Profile)}
	err := s.Set(profile.Profile{Name: "prod", Address: "https://vault.prod", Path: "secret/prod"})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	p, err := s.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if p.Name != "prod" {
		t.Errorf("name = %q, want %q", p.Name, "prod")
	}
}

func TestSet_MissingFields(t *testing.T) {
	s := &profile.Store{Profiles: make(map[string]profile.Profile)}
	if err := s.Set(profile.Profile{Name: "", Address: "http://x", Path: "p"}); err == nil {
		t.Error("expected error for empty name")
	}
	if err := s.Set(profile.Profile{Name: "x", Address: "", Path: "p"}); err == nil {
		t.Error("expected error for empty address")
	}
	if err := s.Set(profile.Profile{Name: "x", Address: "http://x", Path: ""}); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestDelete_Existing(t *testing.T) {
	s := &profile.Store{Profiles: map[string]profile.Profile{
		"staging": {Name: "staging", Address: "http://vault", Path: "secret/staging"},
	}}
	if err := s.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Get("staging"); err == nil {
		t.Error("expected error after delete")
	}
}

func TestDelete_NotFound(t *testing.T) {
	s := &profile.Store{Profiles: make(map[string]profile.Profile)}
	if err := s.Delete("ghost"); err == nil {
		t.Error("expected error deleting non-existent profile")
	}
}

func TestSave_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "profiles.json")
	s := &profile.Store{Profiles: make(map[string]profile.Profile)}
	if err := profile.Save(path, s); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
