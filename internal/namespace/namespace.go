// Package namespace provides utilities for resolving and normalizing
// Vault secret paths across multiple namespaces or environments.
package namespace

import (
	"fmt"
	"strings"
)

// Options configures namespace resolution.
type Options struct {
	// Base is the root path prefix (e.g. "secret/data").
	Base string
	// Env is the environment segment (e.g. "prod", "staging").
	Env string
	// Team is an optional team or service segment.
	Team string
}

// Resolver builds full Vault paths from logical secret names.
type Resolver struct {
	opts Options
}

// New creates a Resolver with the given options.
func New(opts Options) *Resolver {
	return &Resolver{opts: opts}
}

// Resolve returns the full Vault path for the given secret name.
// Path format: <base>/<env>/<team>/<name> (team omitted if empty).
func (r *Resolver) Resolve(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("namespace: secret name must not be empty")
	}
	if r.opts.Base == "" {
		return "", fmt.Errorf("namespace: base path must not be empty")
	}
	if r.opts.Env == "" {
		return "", fmt.Errorf("namespace: env must not be empty")
	}

	parts := []string{r.opts.Base, r.opts.Env}
	if r.opts.Team != "" {
		parts = append(parts, r.opts.Team)
	}
	parts = append(parts, name)

	return strings.Join(parts, "/"), nil
}

// ResolveAll resolves multiple secret names, returning paths in the same order.
func (r *Resolver) ResolveAll(names []string) ([]string, error) {
	paths := make([]string, 0, len(names))
	for _, n := range names {
		p, err := r.Resolve(n)
		if err != nil {
			return nil, err
		}
		paths = append(paths, p)
	}
	return paths, nil
}
