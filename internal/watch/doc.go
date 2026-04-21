// Package watch implements a polling loop that periodically re-fetches secrets
// from HashiCorp Vault and rewrites the local .env file.
//
// Usage:
//
//	w := watch.New(watch.Options{
//		Interval:   30 * time.Second,
//		VaultPath:  "secret/data/myapp",
//		OutputFile: ".env",
//		Out:        os.Stdout,
//	}, fetcher, writer)
//
//	// blocks until context is cancelled or a signal is received
//	watch.RunWithSignals(w, os.Stderr)
package watch
