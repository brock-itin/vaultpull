package env

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// interpolatePattern matches ${VAR} and $VAR style references.
var interpolatePattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateOptions controls interpolation behaviour.
type InterpolateOptions struct {
	// Env is the source map to resolve references from.
	// If nil, os.Getenv is used as fallback.
	Env map[string]string

	// AllowMissing suppresses errors when a referenced variable is not found.
	AllowMissing bool
}

// Interpolate replaces variable references in each value of the provided map
// using values from opts.Env (falling back to the OS environment).
// Returns a new map with substituted values and does not mutate the input.
func Interpolate(secrets map[string]string, opts *InterpolateOptions) (map[string]string, error) {
	if opts == nil {
		opts = &InterpolateOptions{AllowMissing: true}
	}

	result := make(map[string]string, len(secrets))

	for k, v := range secrets {
		expanded, err := interpolateValue(v, secrets, opts)
		if err != nil {
			return nil, fmt.Errorf("interpolate key %q: %w", k, err)
		}
		result[k] = expanded
	}

	return result, nil
}

func interpolateValue(val string, env map[string]string, opts *InterpolateOptions) (string, error) {
	var lastErr error

	replaced := interpolatePattern.ReplaceAllStringFunc(val, func(match string) string {
		name := extractVarName(match)

		if v, ok := env[name]; ok {
			return v
		}
		if v, ok := opts.Env[name]; ok {
			return v
		}
		if v := os.Getenv(name); v != "" {
			return v
		}
		if !opts.AllowMissing {
			lastErr = fmt.Errorf("variable %q not found", name)
			return match
		}
		return match
	})

	if lastErr != nil {
		return "", lastErr
	}
	return replaced, nil
}

func extractVarName(match string) string {
	if strings.HasPrefix(match, "${") {
		return strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")
	}
	return strings.TrimPrefix(match, "$")
}
