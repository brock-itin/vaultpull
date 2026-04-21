package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScan_KeyFound(t *testing.T) {
	input := `DB_HOST=localhost
DB_PORT=5432
`
	r := strings.NewReader(input)
	res, err := Scan(r, "DB_PORT", DefaultScanOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Found {
		t.Fatal("expected key to be found")
	}
	if res.Value != "5432" {
		t.Errorf("expected value 5432, got %q", res.Value)
	}
	if res.Line != 2 {
		t.Errorf("expected line 2, got %d", res.Line)
	}
}

func TestScan_KeyNotFound(t *testing.T) {
	input := `FOO=bar
`
	r := strings.NewReader(input)
	res, err := Scan(r, "MISSING", DefaultScanOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Found {
		t.Fatal("expected key not to be found")
	}
}

func TestScan_CapturesComments(t *testing.T) {
	input := `# database host
# used by app
DB_HOST=localhost
`
	r := strings.NewReader(input)
	res, err := Scan(r, "DB_HOST", DefaultScanOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(res.Comments))
	}
}

func TestScan_CaseInsensitive(t *testing.T) {
	input := `db_host=myhost
`
	r := strings.NewReader(input)
	opts := ScanOptions{CaseSensitive: false}
	res, err := Scan(r, "DB_HOST", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Found {
		t.Fatal("expected case-insensitive match")
	}
}

func TestScan_QuotedValue(t *testing.T) {
	input := `SECRET="my secret value"
`
	r := strings.NewReader(input)
	res, err := Scan(r, "SECRET", DefaultScanOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Value != "my secret value" {
		t.Errorf("expected unquoted value, got %q", res.Value)
	}
}

func TestScanFile_ReturnsErrorOnMissing(t *testing.T) {
	_, err := ScanFile("/nonexistent/.env", "KEY", DefaultScanOptions())
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestScanFile_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	_ = os.WriteFile(path, []byte("API_KEY=abc123\n"), 0o600)

	res, err := ScanFile(path, "API_KEY", DefaultScanOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Found || res.Value != "abc123" {
		t.Errorf("unexpected result: %+v", res)
	}
}
