// Package circuit provides a circuit breaker for wrapping Vault API calls.
//
// The circuit breaker transitions between three states:
//
//   - Closed: normal operation; all calls are forwarded.
//   - Open: too many consecutive failures have occurred; calls are blocked
//     immediately with ErrOpen to avoid hammering a degraded Vault.
//   - Half-Open: after the configured timeout the breaker allows a single
//     probe call through. Success closes the circuit; failure reopens it.
//
// Example:
//
//	b := circuit.New(circuit.DefaultOptions())
//	err := b.Do(func() error {
//		_, err := vaultClient.GetSecrets(path)
//		return err
//	})
//	if errors.Is(err, circuit.ErrOpen) {
//		log.Println("vault unavailable, using cached secrets")
//	}
package circuit
