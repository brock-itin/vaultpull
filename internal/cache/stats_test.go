package cache_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/cache"
)

func newTracked(ttl time.Duration) *cache.TrackedCache {
	return &cache.TrackedCache{Cache: cache.New(ttl)}
}

func TestStats_HitRecorded(t *testing.T) {
	tc := newTracked(time.Minute)
	tc.Cache.Set("secret/x", map[string]string{"K": "v"})
	tc.Get("secret/x")
	if tc.Stats.Hits() != 1 {
		t.Errorf("expected 1 hit, got %d", tc.Stats.Hits())
	}
	if tc.Stats.Misses() != 0 {
		t.Errorf("expected 0 misses, got %d", tc.Stats.Misses())
	}
}

func TestStats_MissRecorded(t *testing.T) {
	tc := newTracked(time.Minute)
	tc.Get("secret/missing")
	if tc.Stats.Misses() != 1 {
		t.Errorf("expected 1 miss, got %d", tc.Stats.Misses())
	}
	if tc.Stats.Hits() != 0 {
		t.Errorf("expected 0 hits, got %d", tc.Stats.Hits())
	}
}

func TestStats_Reset(t *testing.T) {
	tc := newTracked(time.Minute)
	tc.Stats.RecordHit()
	tc.Stats.RecordMiss()
	tc.Stats.Reset()
	if tc.Stats.Hits() != 0 || tc.Stats.Misses() != 0 {
		t.Error("expected counters to be zero after reset")
	}
}

func TestStats_MultipleOps(t *testing.T) {
	tc := newTracked(time.Minute)
	tc.Cache.Set("p", map[string]string{"a": "b"})
	for i := 0; i < 3; i++ {
		tc.Get("p")
	}
	tc.Get("missing")
	if tc.Stats.Hits() != 3 {
		t.Errorf("expected 3 hits, got %d", tc.Stats.Hits())
	}
	if tc.Stats.Misses() != 1 {
		t.Errorf("expected 1 miss, got %d", tc.Stats.Misses())
	}
}
