package expire

import (
	"fmt"
	"io"
	"time"
)

// Print writes a human-readable expiration report to w.
func Print(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No expiration data.")
		return
	}
	for _, r := range results {
		switch r.Status {
		case Expired:
			fmt.Fprintf(w, "[EXPIRED]  %s (expired %s ago)\n", r.Key, formatDuration(-r.TTL))
		case Warning:
			fmt.Fprintf(w, "[WARNING]  %s (expires in %s)\n", r.Key, formatDuration(r.TTL))
		default:
			if r.ExpiresAt.IsZero() {
				fmt.Fprintf(w, "[OK]       %s (no expiry)\n", r.Key)
			} else {
				fmt.Fprintf(w, "[OK]       %s (expires in %s)\n", r.Key, formatDuration(r.TTL))
			}
		}
	}
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	if days > 0 {
		return fmt.Sprintf("%dd%dh", days, hours)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
