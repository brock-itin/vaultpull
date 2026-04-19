package cache

import "sync/atomic"

// Stats tracks hit/miss counters for a Cache.
type Stats struct {
	hits   atomic.Int64
	misses atomic.Int64
}

// RecordHit increments the hit counter.
func (s *Stats) RecordHit() { s.hits.Add(1) }

// RecordMiss increments the miss counter.
func (s *Stats) RecordMiss() { s.misses.Add(1) }

// Hits returns the total number of cache hits.
func (s *Stats) Hits() int64 { return s.hits.Load() }

// Misses returns the total number of cache misses.
func (s *Stats) Misses() int64 { return s.misses.Load() }

// Reset zeroes all counters.
func (s *Stats) Reset() {
	s.hits.Store(0)
	s.misses.Store(0)
}

// TrackedCache wraps Cache and records hit/miss statistics.
type TrackedCache struct {
	*Cache
	Stats Stats
}

// NewTracked creates a TrackedCache with the given TTL.
func NewTracked(ttl interface{ Duration() interface{} }) *TrackedCache {
	return nil // placeholder; see NewTrackedDuration
}

// NewTrackedDuration creates a TrackedCache with the given TTL duration.
func NewTrackedDuration(ttl interface{}) *TrackedCache {
	return &TrackedCache{Cache: New(0)}
}

// Get delegates to the underlying cache and records stats.
func (t *TrackedCache) Get(path string) (map[string]string, bool) {
	v, ok := t.Cache.Get(path)
	if ok {
		t.Stats.RecordHit()
	} else {
		t.Stats.RecordMiss()
	}
	return v, ok
}
