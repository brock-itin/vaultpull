// Package validate checks secrets maps for common issues such as empty values,
// suspicious placeholders, and required key presence.
package validate

import (
	"fmt"
	"strings"
)

// Issue represents a single validation problem.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) Error() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

// Options controls validation behaviour.
type Options struct {
	// Required lists keys that must be present and non-empty.
	Required []string
	// ForbidPlaceholders rejects values that look like unset template placeholders.
	ForbidPlaceholders bool
}

var placeholders = []string{"CHANGEME", "TODO", "FIXME", "<", ">"}

// Check validates secrets against opts and returns all issues found.
func Check(secrets map[string]string, opts Options) []Issue {
	var issues []Issue

	for _, key := range opts.Required {
		v, ok := secrets[key]
		if !ok || strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{Key: key, Message: "required key is missing or empty"})
		}
	}

	if opts.ForbidPlaceholders {
		for k, v := range secrets {
			up := strings.ToUpper(v)
			for _, p := range placeholders {
				if strings.Contains(up, strings.ToUpper(p)) {
					issues = append(issues, Issue{Key: k, Message: fmt.Sprintf("value looks like a placeholder (%s)", p)})
					break
				}
			}
		}
	}

	return issues
}

// HasIssues returns true when Check finds at least one issue.
func HasIssues(secrets map[string]string, opts Options) bool {
	return len(Check(secrets, opts)) > 0
}
