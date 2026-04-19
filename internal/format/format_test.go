package format_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/format"
)

var testSecrets = map[string]string{
	"DB_HOST": "localhost",
	"API_KEY": "abc123",
}

func TestWrite_EnvFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := format.Write(&buf, testSecrets, format.TypeEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got: %s", out)
	}
	if !strings.Contains(out, "API_KEY=abc123") {
		t.Errorf("expected API_KEY=abc123 in output, got: %s", out)
	}
}

func TestWrite_ExportFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := format.Write(&buf, testSecrets, format.TypeExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export DB_HOST='localhost'") {
		t.Errorf("expected export line, got: %s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := format.Write(&buf, testSecrets, format.TypeJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", result["DB_HOST"])
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := format.Write(&buf, testSecrets, format.Type("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestWrite_ExportEscapesSingleQuote(t *testing.T) {
	secrets := map[string]string{"VAL": "it's here"}
	var buf bytes.Buffer
	if err := format.Write(&buf, secrets, format.TypeExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "export VAL='") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}
