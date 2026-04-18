// Package transform provides key transformation utilities for secrets
// pulled from Vault, such as renaming or prefixing keys before writing
// them to .env files.
package transform

import "strings"

// Rule defines a transformation to apply to secret keys.
type Rule struct {
	// AddPrefix prepends a string to every key.
	AddPrefix string
	// StripPrefix removes a leading string from every key.
	StripPrefix string
	// Rename maps original key names to new names.
	Rename map[string]string
}

// Apply applies the transformation rule to the given secrets map and
// returns a new map with transformed keys. Values are unchanged.
func Apply(secrets map[string]string, rule Rule) map[string]string {
	if rule.AddPrefix == "" && rule.StripPrefix == "" && len(rule.Rename) == 0 {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		newKey := k

		if rule.StripPrefix != "" {
			newKey = strings.TrimPrefix(newKey, rule.StripPrefix)
		}

		if rule.AddPrefix != "" {
			newKey = rule.AddPrefix + newKey
		}

		if renamed, ok := rule.Rename[newKey]; ok {
			newKey = renamed
		} else if renamed, ok := rule.Rename[k]; ok {
			newKey = renamed
		}

		result[newKey] = v
	}
	return result
}

// Uppercase converts all keys in the secrets map to uppercase.
func Uppercase(secrets map[string]string) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[strings.ToUpper(k)] = v
	}
	return result
}
