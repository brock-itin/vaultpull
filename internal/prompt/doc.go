// Package prompt provides interactive terminal prompts for vaultpull.
//
// It exposes a Confirmer interface so that callers can swap in AutoConfirm
// when running in non-interactive or CI environments (e.g. via a --yes flag),
// and use the real Prompt implementation when user confirmation is required
// before writing or overwriting secrets in .env files.
package prompt
