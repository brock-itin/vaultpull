// Package filter provides key filtering for secrets pulled from Vault.
package filter

import (
	"strings"
)

// Rule defines how secrets should be filtered before writing.
type Rule struct {
	// Prefix restricts keys to those starting with the given prefix.
	Prefix string
	// Only includes only the listed keys (if non-empty).
	Only []string
	// Exclude removes the listed keys from the result.
	Exclude []string
}

// Apply filters the given secrets map according to the rule and returns
// a new map containing only the keys that pass all conditions.
func Apply(secrets map[string]string, rule Rule) map[string]string {
	only := toSet(rule.Only)
	exclude := toSet(rule.Exclude)

	result := make(map[string]string)
	for k, v := range secrets {
		if rule.Prefix != "" && !strings.HasPrefix(k, rule.Prefix) {
			continue
		}
		if len(only) > 0 && !only[k] {
			continue
		}
		if exclude[k] {
			continue
		}
		result[k] = v
	}
	return result
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
