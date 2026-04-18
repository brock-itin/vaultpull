package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/vaultpull/internal/audit"
)

func TestFileLogger_WritesAndCloses(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit", "vaultpull.log")

	fl, err := audit.NewFileLogger(logPath)
	if err != nil {
		t.Fatalf("NewFileLogger: %v", err)
	}

	if err := fl.Success("secret/test", ".env", []string{"FOO", "BAR"}); err != nil {
		t.Fatalf("Success: %v", err)
	}
	if err := fl.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	f, err := os.Open(logPath)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected at least one log line")
	}

	var entry audit.Entry
	if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if entry.Event != "sync_success" {
		t.Errorf("expected sync_success, got %s", entry.Event)
	}
}

func TestFileLogger_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "a", "b", "c", "audit.log")

	fl, err := audit.NewFileLogger(logPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer fl.Close()

	if _, err := os.Stat(filepath.Dir(logPath)); err != nil {
		t.Errorf("parent dirs not created: %v", err)
	}
}
