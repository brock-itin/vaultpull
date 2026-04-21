package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExport_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env.sh")

	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := Export(out, secrets, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "export BAZ='qux'") {
		t.Errorf("expected BAZ line, got:\n%s", content)
	}
	if !strings.Contains(content, "export FOO='bar'") {
		t.Errorf("expected FOO line, got:\n%s", content)
	}
}

func TestExport_SortedOutput(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env.sh")

	secrets := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"}
	if err := Export(out, secrets, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "export ALPHA") {
		t.Errorf("expected ALPHA first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "export ZEBRA") {
		t.Errorf("expected ZEBRA last, got: %s", lines[2])
	}
}

func TestExport_WithPrefix(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env.sh")

	opts := &ExportOptions{Prefix: "APP_", Overwrite: true, Perm: 0600}
	secrets := map[string]string{"KEY": "value"}
	if err := Export(out, secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "export APP_KEY='value'") {
		t.Errorf("expected prefixed key, got: %s", string(data))
	}
}

func TestExport_NoOverwrite(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env.sh")
	_ = os.WriteFile(out, []byte("existing"), 0600)

	opts := &ExportOptions{Overwrite: false, Perm: 0600}
	err := Export(out, map[string]string{"X": "y"}, opts)
	if err == nil {
		t.Fatal("expected error for no-overwrite on existing file")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestExport_EscapesSingleQuote(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env.sh")

	secrets := map[string]string{"PASS": "it's'fine"}
	if err := Export(out, secrets, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "'it'\\''s'\\''fine'") {
		t.Errorf("expected escaped single quotes, got: %s", string(data))
	}
}

func TestExport_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "nested", "deep", "env.sh")

	if err := Export(out, map[string]string{"K": "v"}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
