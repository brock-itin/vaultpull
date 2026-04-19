// Package template provides Go text/template–based rendering for secret
// injection. It allows users to define a custom .env (or any text) template
// that is populated with secrets fetched from Vault.
//
// Templates receive the secrets map as their data value and may use the
// built-in helper functions:
//
//	{{ required "KEY" . }}  — renders the value or returns an error if absent
//	{{ default "val" "KEY" . }} — renders the value or falls back to a default
//
// Custom delimiters can be set via Options to avoid conflicts with other
// templating systems.
package template
