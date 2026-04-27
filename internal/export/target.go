package export

import (
	"fmt"
	"strings"
)

// Target describes a named export destination with its format and path.
type Target struct {
	Name   string
	Path   string
	Format Format
	Prefix string
}

// ParseTarget parses a target definition string of the form:
//
//	name:path[@format][+prefix]
//
// Examples:
//
//	prod:/etc/app/.env
//	ci:/tmp/ci.env@docker
//	local:.env@shell+APP_
func ParseTarget(s string) (Target, error) {
	if s == "" {
		return Target{}, fmt.Errorf("export: empty target definition")
	}

	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return Target{}, fmt.Errorf("export: invalid target %q: expected name:path", s)
	}

	t := Target{
		Name:   strings.TrimSpace(parts[0]),
		Format: FormatDotEnv,
	}

	rest := parts[1]

	// extract optional prefix after '+'
	if idx := strings.LastIndex(rest, "+"); idx != -1 {
		t.Prefix = rest[idx+1:]
		rest = rest[:idx]
	}

	// extract optional format after '@'
	if idx := strings.LastIndex(rest, "@"); idx != -1 {
		t.Format = Format(strings.TrimSpace(rest[idx+1:]))
		rest = rest[:idx]
	}

	t.Path = strings.TrimSpace(rest)
	if t.Path == "" {
		return Target{}, fmt.Errorf("export: target %q has empty path", s)
	}
	if t.Name == "" {
		return Target{}, fmt.Errorf("export: target %q has empty name", s)
	}

	return t, nil
}

// ExportAll runs Export for each target using the provided secrets.
func ExportAll(secrets map[string]string, targets []Target, baseOpts Options) []error {
	var errs []error
	for _, tgt := range targets {
		opts := baseOpts
		opts.Format = tgt.Format
		opts.OutputPath = tgt.Path
		opts.Prefix = tgt.Prefix
		if err := Export(secrets, opts); err != nil {
			errs = append(errs, fmt.Errorf("target %q: %w", tgt.Name, err))
		}
	}
	return errs
}
