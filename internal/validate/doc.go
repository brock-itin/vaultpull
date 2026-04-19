// Package validate provides pre-write validation for secrets pulled from Vault.
//
// It checks for common problems such as missing required keys, empty values,
// and values that appear to be unfilled placeholders (e.g. CHANGEME, TODO).
//
// Usage:
//
//	issues := validate.Check(secrets, validate.Options{
//		Required:           []string{"DB_PASSWORD", "API_KEY"},
//		ForbidPlaceholders: true,
//	})
package validate
