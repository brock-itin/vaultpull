// Package filter provides utilities for restricting which secret keys
// are written to the local .env file.
//
// A Rule can specify a key prefix, an allowlist (Only), or a denylist
// (Exclude). All conditions are applied together — a key must satisfy
// every non-empty condition to be included in the output.
package filter
