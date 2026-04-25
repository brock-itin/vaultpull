// Package pin implements version pinning for HashiCorp Vault secret paths.
//
// A pin records the specific version of a secret that a project has
// declared stable. Pins are persisted to a JSON file (e.g. .vaultpins)
// alongside the project's other vaultpull configuration.
//
// Typical usage:
//
//	pf, err := pin.Load(".vaultpins")
//	pf.Set("secret/data/myapp", 4, "ci-bot")
//	pin.Save(".vaultpins", pf)
//
// Drift detection compares pinned versions against the currently observed
// versions returned by Vault, surfacing paths that have changed since the
// pin was recorded.
package pin
