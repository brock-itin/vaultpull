// Package rotate provides secret rotation detection for vaultpull.
//
// It evaluates fetched secrets against a configurable age policy
// and identifies which keys are considered stale and should be
// re-fetched or flagged for rotation.
//
// Usage:
//
//	policy := rotate.Policy{MaxAge: 24 * time.Hour}
//	results := rotate.Check(entries, policy, time.Now())
//	if rotate.HasStale(results) {
//		fmt.Println("Stale keys:", rotate.StaleKeys(results))
//	}
package rotate
