// Package sanitize provides utilities for cleaning and normalizing
// secret key names before writing them to .env files.
package sanitize

import (
	"strings"
	"unicode"
)

// Options controls sanitization behavior.
type Options struct {
	// UppercaseKeys converts all key names to uppercase.
	UppercaseKeys bool
	// ReplaceHyphens replaces hyphens in key names with underscores.
	ReplaceHyphens bool
	// StripInvalidChars removes characters not valid in env var names.
	StripInvalidChars bool
}

// DefaultOptions returns sensible sanitization defaults.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:     true,
		ReplaceHyphens:    true,
		StripInvalidChars: true,
	}
}

// Key sanitizes a single env var key name according to the given options.
// A valid env var key contains only letters, digits, and underscores,
// and must not start with a digit.
func Key(k string, opts Options) string {
	if opts.ReplaceHyphens {
		k = strings.ReplaceAll(k, "-", "_")
	}
	if opts.UppercaseKeys {
		k = strings.ToUpper(k)
	}
	if opts.StripInvalidChars {
		k = stripInvalid(k)
	}
	// Strip leading digits to ensure valid identifier.
	k = strings.TrimLeftFunc(k, unicode.IsDigit)
	return k
}

// Map applies Key sanitization to all keys in the provided map,
// returning a new map with sanitized keys. If two keys collide after
// sanitization, the last one (in iteration order) wins.
func Map(secrets map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		sanitized := Key(k, opts)
		if sanitized == "" {
			continue
		}
		out[sanitized] = v
	}
	return out
}

// stripInvalid removes characters that are not letters, digits, or underscores.
func stripInvalid(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
