// Package ttl evaluates time-to-live constraints for secrets pulled from
// HashiCorp Vault.
//
// Each secret path can carry a TTL duration recorded at fetch time. The
// package compares the elapsed time against that duration and classifies
// each entry as Fresh, Warning, or Expired:
//
//   - Fresh   – the secret is well within its TTL.
//   - Warning – less than WarnThreshold (default 20%) of the TTL remains;
//               the caller should schedule a refresh soon.
//   - Expired – the TTL has elapsed; the secret must be re-fetched before
//               use.
//
// Entries with a zero TTL are always considered Fresh (no expiry enforced).
//
// Example:
//
//	results := ttl.Check(entries, ttl.DefaultOptions())
//	if ttl.HasExpired(results) {
//		log.Fatal("one or more secrets have expired")
//	}
package ttl
