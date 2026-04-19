package expire_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/expire"
)

func TestCheck_Fresh(t *testing.T) {
	secrets := map[string]time.Time{
		"API_KEY": time.Now().Add(30 * 24 * time.Hour),
	}
	results := expire.Check(secrets, expire.DefaultOptions())
	if len(results) != 1 || results[0].Status != expire.Fresh {
		t.Errorf("expected Fresh, got %v", results[0].Status)
	}
}

func TestCheck_Warning(t *testing.T) {
	secrets := map[string]time.Time{
		"DB_PASS": time.Now().Add(3 * 24 * time.Hour),
	}
	results := expire.Check(secrets, expire.DefaultOptions())
	if len(results) != 1 || results[0].Status != expire.Warning {
		t.Errorf("expected Warning, got %v", results[0].Status)
	}
}

func TestCheck_Expired(t *testing.T) {
	secrets := map[string]time.Time{
		"OLD_TOKEN": time.Now().Add(-1 * time.Hour),
	}
	results := expire.Check(secrets, expire.DefaultOptions())
	if len(results) != 1 || results[0].Status != expire.Expired {
		t.Errorf("expected Expired, got %v", results[0].Status)
	}
}

func TestCheck_ZeroTime_IsFresh(t *testing.T) {
	secrets := map[string]time.Time{
		"NO_EXPIRY": {},
	}
	results := expire.Check(secrets, expire.DefaultOptions())
	if results[0].Status != expire.Fresh {
		t.Errorf("expected Fresh for zero time")
	}
}

func TestHasExpired(t *testing.T) {
	results := []expire.Result{
		{Key: "A", Status: expire.Fresh},
		{Key: "B", Status: expire.Expired},
	}
	if !expire.HasExpired(results) {
		t.Error("expected HasExpired true")
	}
}

func TestHasWarnings(t *testing.T) {
	results := []expire.Result{
		{Key: "A", Status: expire.Warning},
	}
	if !expire.HasWarnings(results) {
		t.Error("expected HasWarnings true")
	}
}

func TestHasExpired_None(t *testing.T) {
	results := []expire.Result{
		{Key: "A", Status: expire.Fresh},
	}
	if expire.HasExpired(results) {
		t.Error("expected HasExpired false")
	}
}
