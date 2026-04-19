package audit_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/example/vaultpull/internal/audit"
)

func TestLog_Success(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Success("secret/app", ".env", []string{"DB_URL", "API_KEY"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if entry.Event != "sync_success" {
		t.Errorf("expected sync_success, got %s", entry.Event)
	}
	if entry.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", entry.Path)
	}
	if len(entry.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(entry.Keys))
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLog_Failure(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Failure("secret/app", errors.New("forbidden")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if entry.Event != "sync_failure" {
		t.Errorf("expected sync_failure, got %s", entry.Event)
	}
	if !strings.Contains(entry.Error, "forbidden") {
		t.Errorf("expected error to contain 'forbidden', got %s", entry.Error)
	}
}

func TestLog_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	_ = l.Success("secret/a", ".env", []string{"X"})
	_ = l.Success("secret/b", ".env2", []string{"Y"})

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 log lines, got %d", len(lines))
	}
}

func TestLog_SuccessEmptyKeys(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Success("secret/empty", ".env", []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(entry.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(entry.Keys))
	}
}
