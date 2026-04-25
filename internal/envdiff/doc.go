// Package envdiff provides utilities for comparing two env variable maps and
// reporting which keys were added, removed, changed, or left unchanged.
//
// It is used by vaultpull to present a human-readable summary of what will
// change before secrets are written to disk, supporting both plain and
// value-masked output modes.
package envdiff
