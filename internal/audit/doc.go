// Package audit provides structured JSON audit logging for vaultpull.
//
// Each sync operation — successful or failed — is recorded as a JSON line
// containing a timestamp, event type, vault path, output file, affected
// secret keys, and any error message.
//
// Usage:
//
//	l := audit.New(os.Stderr)
//	l.Success("secret/myapp", ".env", []string{"DB_URL"})
package audit
