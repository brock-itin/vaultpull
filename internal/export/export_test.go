package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var testSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "s3cr3t",
	"API_KEY":     "abc123",
}

func TestExport_DotEnvToFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.env")

	opts := DefaultOptions()
	opts.OutputPath = out

	if err := Export(testSecrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "DB_HOST=") {
		t.Errorf("expected DB_HOST in output, got:\n%s", content)
	}
}

func TestExport_ShellFormat(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.sh")

	opts := DefaultOptions()
	opts.Format = FormatShell
	opts.OutputPath = out

	if err := Export(testSecrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "export DB_HOST=") {
		t.Errorf("expected 'export' prefix in shell output")
	}
}

func TestExport_DockerFormat(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "docker.env")

	opts := DefaultOptions()
	opts.Format = FormatDocker
	opts.OutputPath = out

	if err := Export(testSecrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	content := string(data)
	if strings.Contains(content, "export ") {
		t.Errorf("docker format should not contain 'export'")
	}
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected plain KEY=VALUE format")
	}
}

func TestExport_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "secrets.json")

	opts := DefaultOptions()
	opts.Format = FormatJSON
	opts.OutputPath = out

	if err := Export(testSecrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	content := string(data)
	if !strings.HasPrefix(content, "{") {
		t.Errorf("expected JSON object")
	}
	if !strings.Contains(content, `"API_KEY"`) {
		t.Errorf("expected API_KEY in JSON output")
	}
}

func TestExport_WithPrefix(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "prefixed.env")

	opts := DefaultOptions()
	opts.Format = FormatDocker
	opts.Prefix = "APP_"
	opts.OutputPath = out

	if err := Export(map[string]string{"HOST": "localhost"}, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "APP_HOST=localhost") {
		t.Errorf("expected prefixed key in output")
	}
}

func TestExport_OmitEmpty(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "noempty.env")

	opts := DefaultOptions()
	opts.OmitEmpty = true
	opts.OutputPath = out

	secrets := map[string]string{
		"PRESENT": "value",
		"EMPTY":   "",
	}

	if err := Export(secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	content := string(data)
	if strings.Contains(content, "EMPTY") {
		t.Errorf("expected empty key to be omitted")
	}
	if !strings.Contains(content, "PRESENT") {
		t.Errorf("expected PRESENT key in output")
	}
}

func TestExport_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "nested", "deep", "out.env")

	opts := DefaultOptions()
	opts.OutputPath = out

	if err := Export(testSecrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(out); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestExport_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "secret.env")

	opts := DefaultOptions()
	opts.OutputPath = out

	if err := Export(testSecrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
