package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseTarget_Basic(t *testing.T) {
	tgt, err := ParseTarget("prod:/etc/app/.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tgt.Name != "prod" {
		t.Errorf("expected name 'prod', got %q", tgt.Name)
	}
	if tgt.Path != "/etc/app/.env" {
		t.Errorf("expected path '/etc/app/.env', got %q", tgt.Path)
	}
	if tgt.Format != FormatDotEnv {
		t.Errorf("expected default dotenv format, got %q", tgt.Format)
	}
}

func TestParseTarget_WithFormat(t *testing.T) {
	tgt, err := ParseTarget("ci:/tmp/ci.env@docker")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tgt.Format != FormatDocker {
		t.Errorf("expected docker format, got %q", tgt.Format)
	}
}

func TestParseTarget_WithPrefix(t *testing.T) {
	tgt, err := ParseTarget("local:.env@shell+APP_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tgt.Prefix != "APP_" {
		t.Errorf("expected prefix 'APP_', got %q", tgt.Prefix)
	}
	if tgt.Format != FormatShell {
		t.Errorf("expected shell format, got %q", tgt.Format)
	}
}

func TestParseTarget_EmptyString(t *testing.T) {
	_, err := ParseTarget("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestParseTarget_MissingColon(t *testing.T) {
	_, err := ParseTarget("nodestination")
	if err == nil {
		t.Error("expected error for missing colon")
	}
}

func TestParseTarget_EmptyPath(t *testing.T) {
	_, err := ParseTarget("name:")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestExportAll_Success(t *testing.T) {
	dir := t.TempDir()

	targets := []Target{
		{Name: "env", Path: filepath.Join(dir, "out.env"), Format: FormatDotEnv},
		{Name: "docker", Path: filepath.Join(dir, "docker.env"), Format: FormatDocker},
	}

	secrets := map[string]string{"KEY": "value"}
	errs := ExportAll(secrets, targets, DefaultOptions())
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	for _, tgt := range targets {
		if _, err := os.Stat(tgt.Path); err != nil {
			t.Errorf("expected file %s to exist", tgt.Path)
		}
	}
}

func TestExportAll_PartialError(t *testing.T) {
	dir := t.TempDir()

	targets := []Target{
		{Name: "good", Path: filepath.Join(dir, "good.env"), Format: FormatDotEnv},
		{Name: "bad", Path: "/nonexistent/path/that/cannot/be/created/secret.env", Format: FormatDotEnv},
	}

	secrets := map[string]string{"KEY": "value"}
	errs := ExportAll(secrets, targets, DefaultOptions())
	if len(errs) == 0 {
		t.Fatal("expected at least one error")
	}
	if !strings.Contains(errs[0].Error(), "bad") {
		t.Errorf("expected error to mention target name, got: %v", errs[0])
	}
}
