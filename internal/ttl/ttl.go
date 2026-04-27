// Package ttl provides utilities for tracking and enforcing time-to-live
// constraints on secrets fetched from Vault.
package ttl

import (
	"fmt"
	"time"
)

// Status represents the TTL state of a secret.
type Status int

const (
	StatusFresh   Status = iota // within TTL, no action needed
	StatusWarning               // approaching expiry
	StatusExpired               // past TTL, must refresh
)

// Entry holds TTL metadata for a single secret path.
type Entry struct {
	Path      string
	FetchedAt time.Time
	TTL       time.Duration
}

// Options configures TTL evaluation behaviour.
type Options struct {
	// WarnThreshold is the fraction of TTL remaining that triggers a warning.
	// Defaults to 0.2 (warn when less than 20% of TTL remains).
	WarnThreshold float64
	Now           func() time.Time
}

// DefaultOptions returns sensible TTL evaluation defaults.
func DefaultOptions() Options {
	return Options{
		WarnThreshold: 0.2,
		Now:           time.Now,
	}
}

// Result describes the evaluated TTL state of a single entry.
type Result struct {
	Path      string
	Status    Status
	Remaining time.Duration
	Message   string
}

// Check evaluates the TTL status of each entry using the provided options.
func Check(entries []Entry, opts Options) []Result {
	if opts.Now == nil {
		opts.Now = time.Now
	}
	if opts.WarnThreshold <= 0 {
		opts.WarnThreshold = 0.2
	}

	now := opts.Now()
	results := make([]Result, 0, len(entries))

	for _, e := range entries {
		if e.TTL <= 0 {
			results = append(results, Result{
				Path:    e.Path,
				Status:  StatusFresh,
				Message: "no TTL set",
			})
			continue
		}

		elapsed := now.Sub(e.FetchedAt)
		remaining := e.TTL - elapsed

		var status Status
		var msg string

		switch {
		case remaining <= 0:
			status = StatusExpired
			msg = fmt.Sprintf("expired %s ago", (-remaining).Round(time.Second))
		case float64(remaining) < float64(e.TTL)*opts.WarnThreshold:
			status = StatusWarning
			msg = fmt.Sprintf("expires in %s", remaining.Round(time.Second))
		default:
			status = StatusFresh
			msg = fmt.Sprintf("valid for %s", remaining.Round(time.Second))
		}

		results = append(results, Result{
			Path:      e.Path,
			Status:    status,
			Remaining: remaining,
			Message:   msg,
		})
	}

	return results
}

// HasExpired returns true if any result has StatusExpired.
func HasExpired(results []Result) bool {
	for _, r := range results {
		if r.Status == StatusExpired {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any result has StatusWarning or StatusExpired.
func HasWarnings(results []Result) bool {
	for _, r := range results {
		if r.Status == StatusWarning || r.Status == StatusExpired {
			return true
		}
	}
	return false
}
