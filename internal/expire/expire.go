// Package expire provides TTL-based expiration checks for secrets.
package expire

import (
	"time"
)

// Status represents the expiration state of a secret.
type Status int

const (
	Fresh   Status = iota
	Warning        // within warning threshold
	Expired
)

// Result holds the expiration result for a single key.
type Result struct {
	Key       string
	Status    Status
	ExpiresAt time.Time
	TTL       time.Duration
}

// Options configures expiration thresholds.
type Options struct {
	WarningBefore time.Duration // warn if expiring within this window
}

func DefaultOptions() Options {
	return Options{
		WarningBefore: 7 * 24 * time.Hour,
	}
}

// Check evaluates expiration for a map of key -> expiresAt timestamps.
// Keys with zero time are considered non-expiring (Fresh).
func Check(secrets map[string]time.Time, opts Options) []Result {
	now := time.Now()
	results := make([]Result, 0, len(secrets))
	for key, exp := range secrets {
		r := Result{Key: key, ExpiresAt: exp}
		if exp.IsZero() {
			r.Status = Fresh
			results = append(results, r)
			continue
		}
		r.TTL = exp.Sub(now)
		switch {
		case now.After(exp):
			r.Status = Expired
		case r.TTL <= opts.WarningBefore:
			r.Status = Warning
		default:
			r.Status = Fresh
		}
		results = append(results, r)
	}
	return results
}

// HasExpired returns true if any result is Expired.
func HasExpired(results []Result) bool {
	for _, r := range results {
		if r.Status == Expired {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any result is Warning.
func HasWarnings(results []Result) bool {
	for _, r := range results {
		if r.Status == Warning {
			return true
		}
	}
	return false
}
