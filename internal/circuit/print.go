package circuit

import (
	"fmt"
	"io"
)

// StatusLine returns a short human-readable description of the breaker state.
func StatusLine(b *Breaker) string {
	switch b.State() {
	case StateClosed:
		return "circuit: closed (healthy)"
	case StateOpen:
		return "circuit: open (blocking requests)"
	case StateHalfOpen:
		return "circuit: half-open (probing)"
	default:
		return "circuit: unknown"
	}
}

// Print writes a formatted status summary to w.
func Print(w io.Writer, b *Breaker) {
	fmt.Fprintln(w, StatusLine(b))
}
