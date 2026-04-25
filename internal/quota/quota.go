// Package quota tracks and enforces per-path secret fetch limits
// to prevent accidental bulk reads or runaway watch loops.
package quota

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a path exceeds its allowed fetch count.
var ErrQuotaExceeded = errors.New("quota exceeded")

// DefaultOptions returns sensible defaults for quota enforcement.
func DefaultOptions() Options {
	return Options{
		MaxFetches: 100,
		Window:     time.Minute,
	}
}

// Options configures quota behaviour.
type Options struct {
	// MaxFetches is the maximum number of fetches allowed per path within Window.
	MaxFetches int
	// Window is the rolling time window over which fetches are counted.
	Window time.Duration
}

type entry struct {
	count     int
	windowEnd time.Time
}

// Tracker counts fetches per vault path and enforces limits.
type Tracker struct {
	mu      sync.Mutex
	opts    Options
	counts  map[string]*entry
}

// New creates a new Tracker with the given options.
func New(opts Options) *Tracker {
	if opts.MaxFetches <= 0 {
		opts.MaxFetches = DefaultOptions().MaxFetches
	}
	if opts.Window <= 0 {
		opts.Window = DefaultOptions().Window
	}
	return &Tracker{
		opts:   opts,
		counts: make(map[string]*entry),
	}
}

// Record records a fetch attempt for the given path.
// It returns ErrQuotaExceeded if the path has exceeded its allowed fetches
// within the current window.
func (t *Tracker) Record(path string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	e, ok := t.counts[path]
	if !ok || now.After(e.windowEnd) {
		t.counts[path] = &entry{
			count:     1,
			windowEnd: now.Add(t.opts.Window),
		}
		return nil
	}

	e.count++
	if e.count > t.opts.MaxFetches {
		return fmt.Errorf("%w: path %q reached %d fetches in %s",
			ErrQuotaExceeded, path, e.count, t.opts.Window)
	}
	return nil
}

// Usage returns the current fetch count for a path within the active window.
// If the window has expired or the path is unknown, it returns 0.
func (t *Tracker) Usage(path string) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.counts[path]
	if !ok || time.Now().After(e.windowEnd) {
		return 0
	}
	return e.count
}

// Reset clears all recorded counts for all paths.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts = make(map[string]*entry)
}
