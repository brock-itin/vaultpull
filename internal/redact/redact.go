// Package redact provides utilities for scrubbing sensitive keys from secret maps
// before they are logged, printed, or stored in audit trails.
package redact

import "strings"

// Options controls which keys are redacted.
type Options struct {
	// Keywords causes any key containing one of these substrings (case-insensitive)
	// to be redacted.
	Keywords []string
	// Placeholder is substituted for redacted values. Defaults to "[REDACTED]".
	Placeholder string
}

func (o *Options) placeholder() string {
	if o == nil || o.Placeholder == "" {
		return "[REDACTED]"
	}
	return o.Placeholder
}

func (o *Options) keywords() []string {
	if o == nil {
		return defaultKeywords
	}
	if len(o.Keywords) == 0 {
		return defaultKeywords
	}
	return o.Keywords
}

var defaultKeywords = []string{"password", "secret", "token", "key", "credential", "passwd", "private"}

// Map returns a copy of secrets where sensitive values are replaced with the
// configured placeholder.
func Map(secrets map[string]string, opts *Options) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if isSensitive(k, opts.keywords()) {
			out[k] = opts.placeholder()
		} else {
			out[k] = v
		}
	}
	return out
}

// IsSensitive reports whether the given key name is considered sensitive.
func IsSensitive(key string, opts *Options) bool {
	return isSensitive(key, opts.keywords())
}

func isSensitive(key string, keywords []string) bool {
	lower := strings.ToLower(key)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}
