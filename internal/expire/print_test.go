package expire_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/expire"
)

func TestPrint_Empty(t *testing.T) {
	var buf bytes.Buffer
	expire.Print(&buf, nil)
	if !strings.Contains(buf.String(), "No expiration") {
		t.Errorf("expected no expiration message, got: %s", buf.String())
	}
}

func TestPrint_ShowsExpired(t *testing.T) {
	var buf bytes.Buffer
	results := []expire.Result{
		{Key: "OLD_KEY", Status: expire.Expired, TTL: -2 * time.Hour},
	}
	expire.Print(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "[EXPIRED]") || !strings.Contains(out, "OLD_KEY") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestPrint_ShowsWarning(t *testing.T) {
	var buf bytes.Buffer
	results := []expire.Result{
		{Key: "WARN_KEY", Status: expire.Warning, TTL: 48 * time.Hour},
	}
	expire.Print(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "[WARNING]") || !strings.Contains(out, "WARN_KEY") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestPrint_ShowsFresh_NoExpiry(t *testing.T) {
	var buf bytes.Buffer
	results := []expire.Result{
		{Key: "PERM_KEY", Status: expire.Fresh, ExpiresAt: time.Time{}},
	}
	expire.Print(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "[OK]") || !strings.Contains(out, "no expiry") {
		t.Errorf("unexpected output: %s", out)
	}
}
