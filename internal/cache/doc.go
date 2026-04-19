// Package cache implements a lightweight in-memory TTL cache for Vault secret
// paths. It is intended to reduce redundant network calls when the same path
// is resolved multiple times within a single vaultpull invocation.
//
// Entries expire automatically based on the TTL provided at construction time.
// All operations are safe for concurrent use.
package cache
