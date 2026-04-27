package ttl_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/ttl"
)

func makeEntry(path string, fetchedAgo, ttlDur time.Duration) ttl.Entry {
	return ttl.Entry{
		Path:      path,
		FetchedAt: time.Now().Add(-fetchedAgo),
		TTL:       ttlDur,
	}
}

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestCheck_Fresh(t *testing.T) {
	now := time.Now()
	opts := ttl.DefaultOptions()
	opts.Now = fixedNow(now)

	entries := []ttl.Entry{{
		Path:      "secret/app",
		FetchedAt: now.Add(-1 * time.Minute),
		TTL:       10 * time.Minute,
	}}

	results := ttl.Check(entries, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != ttl.StatusFresh {
		t.Errorf("expected Fresh, got %v", results[0].Status)
	}
}

func TestCheck_Warning(t *testing.T) {
	now := time.Now()
	opts := ttl.DefaultOptions()
	opts.Now = fixedNow(now)

	// 90% elapsed of a 10-minute TTL → within warn threshold (20%)
	entries := []ttl.Entry{{
		Path:      "secret/app",
		FetchedAt: now.Add(-9 * time.Minute),
		TTL:       10 * time.Minute,
	}}

	results := ttl.Check(entries, opts)
	if results[0].Status != ttl.StatusWarning {
		t.Errorf("expected Warning, got %v", results[0].Status)
	}
}

func TestCheck_Expired(t *testing.T) {
	now := time.Now()
	opts := ttl.DefaultOptions()
	opts.Now = fixedNow(now)

	entries := []ttl.Entry{{
		Path:      "secret/app",
		FetchedAt: now.Add(-15 * time.Minute),
		TTL:       10 * time.Minute,
	}}

	results := ttl.Check(entries, opts)
	if results[0].Status != ttl.StatusExpired {
		t.Errorf("expected Expired, got %v", results[0].Status)
	}
	if results[0].Remaining >= 0 {
		t.Errorf("expected negative remaining, got %v", results[0].Remaining)
	}
}

func TestCheck_NoTTL_IsFresh(t *testing.T) {
	opts := ttl.DefaultOptions()
	entries := []ttl.Entry{{
		Path:      "secret/app",
		FetchedAt: time.Now().Add(-24 * time.Hour),
		TTL:       0,
	}}

	results := ttl.Check(entries, opts)
	if results[0].Status != ttl.StatusFresh {
		t.Errorf("expected Fresh for zero TTL, got %v", results[0].Status)
	}
}

func TestCheck_Mixed(t *testing.T) {
	now := time.Now()
	opts := ttl.DefaultOptions()
	opts.Now = fixedNow(now)

	entries := []ttl.Entry{
		{Path: "a", FetchedAt: now.Add(-1 * time.Minute), TTL: 10 * time.Minute},
		{Path: "b", FetchedAt: now.Add(-9 * time.Minute), TTL: 10 * time.Minute},
		{Path: "c", FetchedAt: now.Add(-15 * time.Minute), TTL: 10 * time.Minute},
	}

	results := ttl.Check(entries, opts)
	if len(results) != 3 {
		t.Fatalf("expected 3 results")
	}
	if results[0].Status != ttl.StatusFresh {
		t.Errorf("a: expected Fresh")
	}
	if results[1].Status != ttl.StatusWarning {
		t.Errorf("b: expected Warning")
	}
	if results[2].Status != ttl.StatusExpired {
		t.Errorf("c: expected Expired")
	}
}

func TestHasExpired_True(t *testing.T) {
	results := []ttl.Result{{Status: ttl.StatusFresh}, {Status: ttl.StatusExpired}}
	if !ttl.HasExpired(results) {
		t.Error("expected HasExpired to be true")
	}
}

func TestHasExpired_False(t *testing.T) {
	results := []ttl.Result{{Status: ttl.StatusFresh}, {Status: ttl.StatusWarning}}
	if ttl.HasExpired(results) {
		t.Error("expected HasExpired to be false")
	}
}

func TestHasWarnings_IncludesExpired(t *testing.T) {
	results := []ttl.Result{{Status: ttl.StatusExpired}}
	if !ttl.HasWarnings(results) {
		t.Error("expected HasWarnings to be true for expired entry")
	}
}
