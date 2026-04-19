// Package multipath resolves and merges secrets from multiple Vault paths.
package multipath

import "fmt"

// Resolver fetches secrets from a Vault-like source.
type Resolver interface {
	GetSecrets(path string) (map[string]string, error)
}

// Result holds the merged secrets and any per-path errors encountered.
type Result struct {
	Secrets map[string]string
	Errors  map[string]error
}

// HasErrors returns true if any path failed to resolve.
func (r Result) HasErrors() bool {
	return len(r.Errors) > 0
}

// Merge fetches secrets from each path and merges them into a single map.
// Later paths take precedence over earlier ones on key conflicts.
// Errors are collected per-path and returned alongside partial results.
func Merge(resolver Resolver, paths []string) Result {
	result := Result{
		Secrets: make(map[string]string),
		Errors:  make(map[string]error),
	}

	for _, path := range paths {
		if path == "" {
			result.Errors[path] = fmt.Errorf("empty path is not allowed")
			continue
		}

		secrets, err := resolver.GetSecrets(path)
		if err != nil {
			result.Errors[path] = err
			continue
		}

		for k, v := range secrets {
			result.Secrets[k] = v
		}
	}

	return result
}

// MergeStrict is like Merge but returns an error immediately on the first failure.
func MergeStrict(resolver Resolver, paths []string) (map[string]string, error) {
	merged := make(map[string]string)

	for _, path := range paths {
		if path == "" {
			return nil, fmt.Errorf("empty path is not allowed")
		}

		secrets, err := resolver.GetSecrets(path)
		if err != nil {
			return nil, fmt.Errorf("path %q: %w", path, err)
		}

		for k, v := range secrets {
			merged[k] = v
		}
	}

	return merged, nil
}
