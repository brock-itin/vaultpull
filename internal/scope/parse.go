package scope

import (
	"fmt"
	"strings"
)

// ParseScopes parses a slice of raw scope strings into Scope values.
// Each entry must follow the format: name:path or name:path:prefix
//
// Examples:
//
//	"app:secret/myapp"
//	"db:secret/database:DB_"
func ParseScopes(raw []string) ([]Scope, error) {
	scopes := make([]Scope, 0, len(raw))
	for _, r := range raw {
		s, err := parseOne(r)
		if err != nil {
			return nil, err
		}
		scopes = append(scopes, s)
	}
	return scopes, nil
}

func parseOne(raw string) (Scope, error) {
	parts := strings.SplitN(raw, ":", 3)
	if len(parts) < 2 {
		return Scope{}, fmt.Errorf("invalid scope %q: expected name:path[:prefix]", raw)
	}
	name := strings.TrimSpace(parts[0])
	path := strings.TrimSpace(parts[1])
	if name == "" {
		return Scope{}, fmt.Errorf("invalid scope %q: name must not be empty", raw)
	}
	if path == "" {
		return Scope{}, fmt.Errorf("invalid scope %q: path must not be empty", raw)
	}
	var prefix string
	if len(parts) == 3 {
		prefix = strings.TrimSpace(parts[2])
	}
	return Scope{Name: name, Path: path, Prefix: prefix}, nil
}
