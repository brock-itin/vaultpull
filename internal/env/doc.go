// Package env handles writing and merging secrets into local .env files.
//
// It supports:
//   - Creating new .env files with secrets fetched from Vault
//   - Preserving existing keys when overwrite is disabled
//   - Optionally backing up the existing .env file before writing
//
// File permissions are set to 0600 to prevent unintended access.
package env
