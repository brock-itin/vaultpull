// Package rotate provides utilities for detecting stale secrets
// and prompting rotation based on age or policy.
package rotate

import (
	"time"
)

// Entry represents a secret with metadata for rotation checks.
type Entry struct {
	Key       string
	FetchedAt time.Time
}

// Policy defines rotation rules.
type Policy struct {
	// MaxAge is the maximum allowed age of a secret before it is considered stale.
	MaxAge time.Duration
}

// Result holds the outcome of a rotation check for a single entry.
type Result struct {
	Key   string
	Stale bool
	Age   time.Duration
}

// Check evaluates each entry against the policy and returns results.
func Check(entries []Entry, policy Policy, now time.Time) []Result {
	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		age := now.Sub(e.FetchedAt)
		results = append(results, Result{
			Key:   e.Key,
			Stale: age > policy.MaxAge,
			Age:   age,
		})
	}
	return results
}

// StaleKeys returns only the keys that are considered stale.
func StaleKeys(results []Result) []string {
	var keys []string
	for _, r := range results {
		if r.Stale {
			keys = append(keys, r.Key)
		}
	}
	return keys
}

// HasStale returns true if any result is stale.
func HasStale(results []Result) bool {
	for _, r := range results {
		if r.Stale {
			return true
		}
	}
	return false
}
