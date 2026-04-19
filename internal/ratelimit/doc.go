// Package ratelimit implements a token-bucket rate limiter used to throttle
// outbound requests to HashiCorp Vault. This prevents overwhelming the Vault
// server when syncing large numbers of secrets or paths in a single run.
//
// Usage:
//
//	l, err := ratelimit.New(ratelimit.DefaultOptions())
//	if err != nil { ... }
//	if err := l.Wait(ctx); err != nil { ... }
package ratelimit
