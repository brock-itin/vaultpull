// Package mask provides utilities for redacting secret values
// before they are displayed in logs or terminal output.
package mask

import "strings"

const defaultMask = "****"

// Options controls masking behaviour.
type Options struct {
	// ShowPrefix reveals the first N characters before masking.
	ShowPrefix int
}

// Value masks a secret string according to opts.
// If opts is nil, the entire value is replaced with ****.
func Value(secret string, opts *Options) string {
	if secret == "" {
		return ""
	}
	if opts == nil || opts.ShowPrefix <= 0 {
		return defaultMask
	}
	n := opts.ShowPrefix
	if n >= len(secret) {
		return defaultMask
	}
	return secret[:n] + strings.Repeat("*", 4)
}

// Map masks every value in a map of secrets.
func Map(secrets map[string]string, opts *Options) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = Value(v, opts)
	}
	return out
}
