// Package ratelimit provides a simple token-bucket rate limiter for Vault API calls.
package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Options configures the rate limiter.
type Options struct {
	// RequestsPerSecond is the maximum number of requests allowed per second.
	RequestsPerSecond int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		RequestsPerSecond: 10,
	}
}

// Limiter controls request rate using a token bucket.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per nanosecond
	lastTick time.Time
}

// New creates a Limiter with the given options.
func New(opts Options) (*Limiter, error) {
	if opts.RequestsPerSecond <= 0 {
		return nil, fmt.Errorf("ratelimit: RequestsPerSecond must be > 0, got %d", opts.RequestsPerSecond)
	}
	max := float64(opts.RequestsPerSecond)
	return &Limiter{
		tokens:   max,
		max:      max,
		rate:     max / float64(time.Second),
		lastTick: time.Now(),
	}, nil
}

// Wait blocks until a token is available or ctx is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		l.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(l.lastTick)
		l.tokens += float64(elapsed) * l.rate
		if l.tokens > l.max {
			l.tokens = l.max
		}
		l.lastTick = now
		if l.tokens >= 1 {
			l.tokens--
			l.mu.Unlock()
			return nil
		}
		waitFor := time.Duration((1-l.tokens)/l.rate)
		l.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitFor):
		}
	}
}
