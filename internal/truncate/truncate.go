// Package truncate provides utilities for truncating secret values
// to safe lengths before displaying or logging them.
package truncate

const defaultMaxLen = 64

// Options configures truncation behaviour.
type Options struct {
	// MaxLen is the maximum number of characters to keep. Defaults to 64.
	MaxLen int
	// Suffix is appended when a value is truncated. Defaults to "...".
	Suffix string
}

func defaults(o *Options) *Options {
	if o == nil {
		o = &Options{}
	}
	if o.MaxLen <= 0 {
		o.MaxLen = defaultMaxLen
	}
	if o.Suffix == "" {
		o.Suffix = "..."
	}
	return o
}

// Value truncates a single string value according to opts.
// If opts is nil, defaults are used.
func Value(s string, opts *Options) string {
	opts = defaults(opts)
	if len(s) <= opts.MaxLen {
		return s
	}
	return s[:opts.MaxLen] + opts.Suffix
}

// Map truncates every value in the provided map, returning a new map.
// Keys are never modified.
func Map(m map[string]string, opts *Options) map[string]string {
	opts = defaults(opts)
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = Value(v, opts)
	}
	return out
}
