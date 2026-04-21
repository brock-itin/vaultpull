package checkpoint_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/checkpoint"
)

func TestPrint_Empty(t *testing.T) {
	var buf bytes.Buffer
	checkpoint.Print(&buf, checkpoint.New(), nil)
	if !strings.Contains(buf.String(), "no checkpoint") {
		t.Errorf("expected empty message, got: %q", buf.String())
	}
}

func TestPrint_ShowsOK(t *testing.T) {
	c := checkpoint.New()
	c.Record(checkpoint.Entry{
		Path:     "secret/app",
		SyncedAt: time.Now(),
		KeyCount: 4,
		Success:  true,
	})
	var buf bytes.Buffer
	checkpoint.Print(&buf, c, nil)
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Errorf("expected [OK] in output: %q", out)
	}
	if !strings.Contains(out, "secret/app") {
		t.Errorf("expected path in output: %q", out)
	}
}

func TestPrint_ShowsFailed(t *testing.T) {
	c := checkpoint.New()
	c.Record(checkpoint.Entry{
		Path:    "secret/broken",
		Success: false,
		Error:   "permission denied",
	})
	var buf bytes.Buffer
	checkpoint.Print(&buf, c, nil)
	out := buf.String()
	if !strings.Contains(out, "[FAILED]") {
		t.Errorf("expected [FAILED] in output: %q", out)
	}
	if !strings.Contains(out, "permission denied") {
		t.Errorf("expected error text in output: %q", out)
	}
}

func TestPrint_ShowsStale(t *testing.T) {
	c := checkpoint.New()
	c.Record(checkpoint.Entry{
		Path:     "secret/old",
		SyncedAt: time.Now().Add(-2 * time.Hour),
		KeyCount: 2,
		Success:  true,
	})
	var buf bytes.Buffer
	opts := &checkpoint.PrintOptions{StaleAfter: time.Hour}
	checkpoint.Print(&buf, c, opts)
	out := buf.String()
	if !strings.Contains(out, "[STALE]") {
		t.Errorf("expected [STALE] in output: %q", out)
	}
}

func TestPrint_SortedPaths(t *testing.T) {
	c := checkpoint.New()
	for _, p := range []string{"secret/z", "secret/a", "secret/m"} {
		c.Record(checkpoint.Entry{Path: p, Success: true, SyncedAt: time.Now()})
	}
	var buf bytes.Buffer
	checkpoint.Print(&buf, c, nil)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "secret/a") {
		t.Errorf("expected first line to be secret/a, got: %q", lines[0])
	}
	if !strings.Contains(lines[2], "secret/z") {
		t.Errorf("expected last line to be secret/z, got: %q", lines[2])
	}
}
