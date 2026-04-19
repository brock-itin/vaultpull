// Package multipath provides utilities for fetching and merging secrets
// from multiple Vault paths into a single key-value map.
//
// Use Merge for a fault-tolerant approach that collects per-path errors
// while returning partial results. Use MergeStrict when all paths must
// succeed for the operation to be considered valid.
package multipath
