package env

import (
	"fmt"
	"os"
)

// LoadOptions controls how environment variables are loaded from a file.
type LoadOptions struct {
	// Overwrite controls whether existing OS env vars are overwritten.
	Overwrite bool
	// Prefix is an optional prefix to strip from keys before setting them.
	Prefix string
}

// DefaultLoadOptions returns sensible defaults for LoadOptions.
func DefaultLoadOptions() LoadOptions {
	return LoadOptions{
		Overwrite: false,
		Prefix:    "",
	}
}

// Load reads a .env file at path and sets the key/value pairs as OS environment
// variables. If opts.Overwrite is false, existing variables are not replaced.
// If opts.Prefix is set, it is stripped from each key before calling os.Setenv.
func Load(path string, opts LoadOptions) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("env: open %q: %w", path, err)
	}
	defer f.Close()

	entries, err := Parse(f)
	if err != nil {
		return fmt.Errorf("env: parse %q: %w", path, err)
	}

	for _, e := range entries {
		key := e.Key
		if opts.Prefix != "" && len(key) > len(opts.Prefix) && key[:len(opts.Prefix)] == opts.Prefix {
			key = key[len(opts.Prefix):]
		}
		if key == "" {
			continue
		}
		if !opts.Overwrite {
			if _, exists := os.LookupEnv(key); exists {
				continue
			}
		}
		if err := os.Setenv(key, e.Value); err != nil {
			return fmt.Errorf("env: setenv %q: %w", key, err)
		}
	}
	return nil
}
