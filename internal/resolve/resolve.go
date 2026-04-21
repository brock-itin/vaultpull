// Package resolve provides utilities for resolving secret keys
// from Vault paths into structured maps, applying optional transformations.
package resolve

import (
	"fmt"
	"strings"
)

// Options controls how resolution behaves.
type Options struct {
	// StripPathPrefix removes the Vault path prefix from resolved keys.
	StripPathPrefix bool
	// Separator replaces path separators in keys. Defaults to "_".
	Separator string
}

// DefaultOptions returns sensible defaults for resolution.
func DefaultOptions() Options {
	return Options{
		StripPathPrefix: true,
		Separator:       "_",
	}
}

// Result holds a resolved key-value pair with its source path.
type Result struct {
	Key    string
	Value  string
	Source string
}

// Resolve takes a map of raw Vault secrets keyed by path and resolves
// them into a flat key-value map using the provided options.
func Resolve(secrets map[string]map[string]string, opts Options) ([]Result, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	var results []Result

	for path, kvs := range secrets {
		for k, v := range kvs {
			resolved, err := resolveKey(path, k, opts)
			if err != nil {
				return nil, fmt.Errorf("resolve: path %q key %q: %w", path, k, err)
			}
			results = append(results, Result{
				Key:    resolved,
				Value:  v,
				Source: path,
			})
		}
	}

	return results, nil
}

// ToMap converts a slice of Results into a plain string map.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Value
	}
	return m
}

func resolveKey(path, key string, opts Options) (string, error) {
	if key == "" {
		return "", fmt.Errorf("empty key")
	}

	if !opts.StripPathPrefix {
		prefix := strings.ReplaceAll(strings.Trim(path, "/"), "/", opts.Separator)
		if prefix != "" {
			return prefix + opts.Separator + key, nil
		}
	}

	return key, nil
}
