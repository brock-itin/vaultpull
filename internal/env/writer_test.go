package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := Write(path, secrets, WriteOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", content)
	}
	if !strings.Contains(content, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got: %s", content)
	}
}

func TestWrite_PreservesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("EXISTING=value\n"), 0o600)

	if err := Write(path, map[string]string{"NEW": "entry"}, WriteOptions{Overwrite: false}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "EXISTING=value") {
		t.Errorf("expected EXISTING=value to be preserved, got: %s", content)
	}
	if !strings.Contains(content, "NEW=entry") {
		t.Errorf("expected NEW=entry in output, got: %s", content)
	}
}

func TestWrite_OverwriteExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("KEY=old\n"), 0o600)

	if err := Write(path, map[string]string{"KEY": "new"}, WriteOptions{Overwrite: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "KEY=new") {
		t.Errorf("expected KEY=new, got: %s", string(data))
	}
}

func TestWrite_BackupCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("ORIG=1\n"), 0o600)

	if err := Write(path, map[string]string{"X": "y"}, WriteOptions{BackupExisting: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path + ".bak"); os.IsNotExist(err) {
		t.Error("expected .env.bak to be created")
	}
}

func TestWrite_NoOverwriteKeepsOldValue(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("KEY=original\n"), 0o600)

	_ = Write(path, map[string]string{"KEY": "ignored"}, WriteOptions{Overwrite: false})

	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "KEY=ignored") {
		t.Errorf("expected KEY=original to be kept, got: %s", string(data))
	}
}
