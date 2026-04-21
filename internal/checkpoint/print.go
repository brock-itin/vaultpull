package checkpoint

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// PrintOptions controls display of checkpoint status.
type PrintOptions struct {
	// StaleAfter marks an entry as stale if older than this duration.
	// Zero disables staleness checks.
	StaleAfter time.Duration
}

// Print writes a human-readable summary of checkpoint entries to w.
func Print(w io.Writer, c *Checkpoint, opts *PrintOptions) {
	if opts == nil {
		opts = &PrintOptions{}
	}
	if len(c.Entries) == 0 {
		fmt.Fprintln(w, "no checkpoint entries recorded")
		return
	}

	paths := make([]string, 0, len(c.Entries))
	for p := range c.Entries {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, p := range paths {
		e := c.Entries[p]
		status := statusLabel(e, opts)
		age := time.Since(e.SyncedAt).Truncate(time.Second)
		fmt.Fprintf(w, "  %-40s  keys=%-3d  age=%-12s  %s\n",
			e.Path, e.KeyCount, age.String(), status)
		if !e.Success && e.Error != "" {
			fmt.Fprintf(w, "    error: %s\n", e.Error)
		}
	}
}

func statusLabel(e Entry, opts *PrintOptions) string {
	if !e.Success {
		return "[FAILED]"
	}
	if opts.StaleAfter > 0 && time.Since(e.SyncedAt) > opts.StaleAfter {
		return "[STALE]"
	}
	return "[OK]"
}
