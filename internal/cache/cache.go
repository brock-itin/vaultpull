// Package cache provides a simple time-based local cache for Vault secrets
// to reduce redundant API calls during a session.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached secret map and its expiry time.
type Entry struct {
	Secrets   map[string]string
	FetchedAt time.Time
	TTL       time.Duration
}

// IsExpired reports whether the cache entry has expired.
func (e Entry) IsExpired() bool {
	if e.TTL <= 0 {
		return true
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// Cache is a thread-safe in-memory store keyed by Vault path.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New creates a Cache with the given default TTL.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Get returns the secrets for a path if present and not expired.
func (c *Cache) Get(path string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[path]
	if !ok || e.IsExpired() {
		return nil, false
	}
	copy := make(map[string]string, len(e.Secrets))
	for k, v := range e.Secrets {
		copy[k] = v
	}
	return copy, true
}

// Set stores secrets for a path using the cache's default TTL.
func (c *Cache) Set(path string, secrets map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	c.entries[path] = Entry{
		Secrets:   copy,
		FetchedAt: time.Now(),
		TTL:       c.ttl,
	}
}

// Invalidate removes a single path from the cache.
func (c *Cache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]Entry)
}
