// Package vault provides a lightweight client for reading secrets
// from a HashiCorp Vault KV v2 engine.
//
// Usage:
//
//	client, err := vault.NewClient(address, token)
//	if err != nil { ... }
//
//	secrets, err := client.GetSecrets("secret/data/myapp")
//	if err != nil { ... }
//
// The returned map contains plain key/value string pairs ready
// to be written to a .env file by the envwriter package.
package vault
