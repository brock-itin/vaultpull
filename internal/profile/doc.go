// Package profile provides named configuration profiles for vaultpull.
//
// A profile captures the Vault address, secret path, and optional output file
// for a named environment (e.g. "dev", "staging", "prod"). Profiles are
// persisted as JSON and can be loaded at runtime to avoid repeating flags.
//
// Example usage:
//
//	store, _ := profile.Load("~/.vaultpull/profiles.json")
//	p, _ := store.Get("dev")
//	fmt.Println(p.Address)
package profile
