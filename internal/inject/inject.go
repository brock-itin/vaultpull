// Package inject provides utilities for injecting secrets into
// process environments and command execution contexts.
package inject

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Options controls how secrets are injected into the environment.
type Options struct {
	// Overwrite replaces existing environment variables if true.
	Overwrite bool
	// Prefix is prepended to every injected key.
	Prefix string
}

// DefaultOptions returns sensible defaults for injection.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
		Prefix:    "",
	}
}

// IntoProcess sets secrets as environment variables in the current process.
// Existing variables are preserved unless opts.Overwrite is true.
func IntoProcess(secrets map[string]string, opts Options) error {
	for k, v := range secrets {
		key := buildKey(k, opts.Prefix)
		if key == "" {
			return fmt.Errorf("inject: empty key derived from %q", k)
		}
		if _, exists := os.LookupEnv(key); exists && !opts.Overwrite {
			continue
		}
		if err := os.Setenv(key, v); err != nil {
			return fmt.Errorf("inject: setenv %q: %w", key, err)
		}
	}
	return nil
}

// IntoCommand returns a copy of cmd with secrets merged into its environment.
// The base environment is inherited from the current process.
func IntoCommand(cmd *exec.Cmd, secrets map[string]string, opts Options) error {
	base := os.Environ()
	existing := make(map[string]bool, len(base))
	for _, e := range base {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			existing[parts[0]] = true
		}
	}

	env := make([]string, 0, len(base)+len(secrets))
	env = append(env, base...)

	for k, v := range secrets {
		key := buildKey(k, opts.Prefix)
		if key == "" {
			return fmt.Errorf("inject: empty key derived from %q", k)
		}
		if existing[key] && !opts.Overwrite {
			continue
		}
		env = append(env, key+"="+v)
	}

	cmd.Env = env
	return nil
}

func buildKey(key, prefix string) string {
	if prefix == "" {
		return key
	}
	return prefix + key
}
