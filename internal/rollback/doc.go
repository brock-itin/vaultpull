// Package rollback restores a local .env file to a previously captured
// snapshot state. It is intended to be used when a sync operation produces
// an undesired result and the operator needs to revert to a known-good
// configuration quickly.
//
// Usage:
//
//	opts := rollback.DefaultOptions()
//	result := rollback.Execute(".env", opts)
//	if result.Err != nil {
//	    log.Fatal(result.Err)
//	}
//	fmt.Printf("Restored %d keys from snapshot taken at %s\n",
//	    result.KeysRestored, result.SnapshotAt)
package rollback
