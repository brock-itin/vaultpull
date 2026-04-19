package rotate_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/rotate"
)

var baseTime = time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)

func makeEntry(key string, ageHours float64) rotate.Entry {
	return rotate.Entry{
		Key:       key,
		FetchedAt: baseTime.Add(-time.Duration(ageHours * float64(time.Hour))),
	}
}

func TestCheck_AllFresh(t *testing.T) {
	entries := []rotate.Entry{makeEntry("DB_PASS", 1), makeEntry("API_KEY", 2)}
	policy := rotate.Policy{MaxAge: 24 * time.Hour}
	results := rotate.Check(entries, policy, baseTime)
	for _, r := range results {
		if r.Stale {
			t.Errorf("expected %s to be fresh", r.Key)
		}
	}
}

func TestCheck_AllStale(t *testing.T) {
	entries := []rotate.Entry{makeEntry("DB_PASS", 48), makeEntry("API_KEY", 72)}
	policy := rotate.Policy{MaxAge: 24 * time.Hour}
	results := rotate.Check(entries, policy, baseTime)
	for _, r := range results {
		if !r.Stale {
			t.Errorf("expected %s to be stale", r.Key)
		}
	}
}

func TestCheck_Mixed(t *testing.T) {
	entries := []rotate.Entry{makeEntry("FRESH", 1), makeEntry("OLD", 50)}
	policy := rotate.Policy{MaxAge: 24 * time.Hour}
	results := rotate.Check(entries, policy, baseTime)
	if results[0].Stale {
		t.Error("FRESH should not be stale")
	}
	if !results[1].Stale {
		t.Error("OLD should be stale")
	}
}

func TestStaleKeys(t *testing.T) {
	results := []rotate.Result{
		{Key: "A", Stale: false},
		{Key: "B", Stale: true},
		{Key: "C", Stale: true},
	}
	keys := rotate.StaleKeys(results)
	if len(keys) != 2 {
		t.Fatalf("expected 2 stale keys, got %d", len(keys))
	}
}

func TestHasStale_True(t *testing.T) {
	results := []rotate.Result{{Key: "X", Stale: true}}
	if !rotate.HasStale(results) {
		t.Error("expected HasStale to return true")
	}
}

func TestHasStale_False(t *testing.T) {
	results := []rotate.Result{{Key: "X", Stale: false}}
	if rotate.HasStale(results) {
		t.Error("expected HasStale to return false")
	}
}
