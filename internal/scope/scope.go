// Package scope provides path-scoped secret resolution,
// allowing secrets to be fetched from multiple Vault paths
// and merged under a common key namespace.
package scope

import (
	"fmt"
	"strings"
)

// Scope represents a named Vault path scope with an optional key prefix.
type Scope struct {
	Name   string
	Path   string
	Prefix string
}

// Options controls how scopes are resolved.
type Options struct {
	// StripPrefix removes the scope prefix from resolved keys.
	StripPrefix bool
}

// DefaultOptions returns sensible defaults for scope resolution.
func DefaultOptions() Options {
	return Options{
		StripPrefix: false,
	}
}

// Resolve applies the scope's prefix to each key in secrets,
// returning a new map with prefixed (or stripped) keys.
func Resolve(s Scope, secrets map[string]string, opts Options) (map[string]string, error) {
	if s.Path == "" {
		return nil, fmt.Errorf("scope %q: path must not be empty", s.Name)
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		key := applyPrefix(k, s.Prefix, opts.StripPrefix)
		out[key] = v
	}
	return out, nil
}

// ResolveAll resolves multiple scopes and merges the results.
// Later scopes win on key conflicts.
func ResolveAll(scopes []Scope, fetch func(path string) (map[string]string, error), opts Options) (map[string]string, error) {
	merged := make(map[string]string)
	for _, s := range scopes {
		secrets, err := fetch(s.Path)
		if err != nil {
			return nil, fmt.Errorf("scope %q (path %q): %w", s.Name, s.Path, err)
		}
		resolved, err := Resolve(s, secrets, opts)
		if err != nil {
			return nil, err
		}
		for k, v := range resolved {
			merged[k] = v
		}
	}
	return merged, nil
}

func applyPrefix(key, prefix string, strip bool) string {
	if prefix == "" {
		return key
	}
	if strip {
		return strings.TrimPrefix(key, prefix)
	}
	return prefix + key
}
