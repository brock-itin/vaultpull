// Package label provides utilities for attaching and filtering key-value
// metadata labels to secrets, enabling grouping and selective sync operations.
package label

import "fmt"

// Labels is a map of string key-value pairs attached to a secret set.
type Labels map[string]string

// Options controls how label matching is performed.
type Options struct {
	// Required is a set of key=value pairs that must all match.
	Required Labels
	// Forbidden is a set of keys whose presence causes rejection.
	Forbidden []string
}

// DefaultOptions returns Options with no constraints.
func DefaultOptions() Options {
	return Options{}
}

// Match reports whether the given labels satisfy the options.
// All Required entries must be present with matching values.
// None of the Forbidden keys may be present.
func Match(labels Labels, opts Options) bool {
	for k, v := range opts.Required {
		got, ok := labels[k]
		if !ok || got != v {
			return false
		}
	}
	for _, k := range opts.Forbidden {
		if _, ok := labels[k]; ok {
			return false
		}
	}
	return true
}

// Parse converts a slice of "key=value" strings into a Labels map.
// It returns an error if any entry is malformed.
func Parse(pairs []string) (Labels, error) {
	out := make(Labels, len(pairs))
	for _, p := range pairs {
		for i := 0; i < len(p); i++ {
			if p[i] == '=' {
				k, v := p[:i], p[i+1:]
				if k == "" {
					return nil, fmt.Errorf("label %q has empty key", p)
				}
				out[k] = v
				goto next
			}
		}
		return nil, fmt.Errorf("label %q missing '='", p)
	next:
	}
	return out, nil
}

// Filter returns only the entries from secrets whose keys have labels
// satisfying opts. secretLabels maps each secret key to its Labels.
func Filter(secrets map[string]string, secretLabels map[string]Labels, opts Options) map[string]string {
	out := make(map[string]string)
	for k, v := range secrets {
		lbls := secretLabels[k]
		if Match(lbls, opts) {
			out[k] = v
		}
	}
	return out
}
