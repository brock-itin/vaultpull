// Package snapshot provides save/load functionality for capturing a point-in-time
// view of secrets pulled from Vault.
//
// Snapshots are stored as JSON files on disk with restricted permissions (0600)
// and can be used to detect drift, enable rollback comparisons, or audit
// what was written during a previous sync.
package snapshot
