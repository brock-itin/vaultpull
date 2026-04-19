// Package retry provides configurable retry logic with exponential backoff
// for use when communicating with Vault or other external services that may
// experience transient failures.
//
// Basic usage:
//
//	err := retry.Do(retry.DefaultOptions(), func(attempt int) error {
//		return callVault()
//	})
package retry
