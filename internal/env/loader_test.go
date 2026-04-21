package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeEnvFile: %v", err)
	}
	return p
}

func TestLoad_SetsVariables(t *testing.T) {
	p := writeEnvFile(t, "LOADER_FOO=bar\nLOADER_BAZ=qux\n")
	t.Setenv("LOADER_FOO", "")
	t.Setenv("LOADER_BAZ", "")

	if err := Load(p, DefaultLoadOptions()); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := os.Getenv("LOADER_FOO"); got != "bar" {
		t.Errorf("LOADER_FOO = %q, want %q", got, "bar")
	}
	if got := os.Getenv("LOADER_BAZ"); got != "qux" {
		t.Errorf("LOADER_BAZ = %q, want %q", got, "qux")
	}
}

func TestLoad_NoOverwrite(t *testing.T) {
	p := writeEnvFile(t, "LOADER_KEEP=new\n")
	t.Setenv("LOADER_KEEP", "original")

	opts := DefaultLoadOptions()
	opts.Overwrite = false
	if err := Load(p, opts); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := os.Getenv("LOADER_KEEP"); got != "original" {
		t.Errorf("LOADER_KEEP = %q, want %q", got, "original")
	}
}

func TestLoad_Overwrite(t *testing.T) {
	p := writeEnvFile(t, "LOADER_OVER=new\n")
	t.Setenv("LOADER_OVER", "original")

	opts := DefaultLoadOptions()
	opts.Overwrite = true
	if err := Load(p, opts); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := os.Getenv("LOADER_OVER"); got != "new" {
		t.Errorf("LOADER_OVER = %q, want %q", got, "new")
	}
}

func TestLoad_StripPrefix(t *testing.T) {
	p := writeEnvFile(t, "APP_LOADER_DB=postgres\n")
	t.Setenv("LOADER_DB", "")

	opts := DefaultLoadOptions()
	opts.Prefix = "APP_"
	opts.Overwrite = true
	if err := Load(p, opts); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got := os.Getenv("LOADER_DB"); got != "postgres" {
		t.Errorf("LOADER_DB = %q, want %q", got, "postgres")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	err := Load("/nonexistent/.env", DefaultLoadOptions())
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
