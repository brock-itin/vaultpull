// Package namespace resolves logical secret names into full Vault paths
// using a configurable base path, environment, and optional team segment.
//
// Example:
//
//	r := namespace.New(namespace.Options{
//		Base: "secret/data",
//		Env:  "prod",
//		Team: "platform",
//	})
//	path, err := r.Resolve("database")
//	// path == "secret/data/prod/platform/database"
package namespace
