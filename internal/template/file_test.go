package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenderFile_Success(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "tmpl.env")
	dst := filepath.Join(dir, "out", ".env")

	content := `DB={{ index . "DB" }}`
	if err := os.WriteFile(src, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	secrets := map[string]string{"DB": "postgres"}
	if err := RenderFile(src, dst, secrets, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "DB=postgres" {
		t.Errorf("got %q", string(got))
	}
}

func TestRenderFile_MissingSrc(t *testing.T) {
	dir := t.TempDir()
	err := RenderFile(filepath.Join(dir, "nope.tmpl"), filepath.Join(dir, "out"), nil, nil)
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestRenderFile_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "t.env")
	dst := filepath.Join(dir, "a", "b", "c", ".env")

	if err := os.WriteFile(src, []byte("X=1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RenderFile(src, dst, map[string]string{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	 err := os.Stat(dst); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

func TestRenderFile_OutputPermissions(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "t.env")
	dst := filepath.Join(dir, ".env")

	if err := os.WriteFile(src, []byte("A=1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RenderFile(src, dst, map[string]string{}, nil); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
Perm() != 0o600 {
		t.Errorf("expected 0600, got %v", info.Mode().Perm())
	}
}
