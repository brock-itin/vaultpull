package hook

import (
	"fmt"
	"io"
)

// Print writes a human-readable summary of hook results to w.
func Print(w io.Writer, results []Result) {
	if len(results) == 0 {
		return
	}
	for _, r := range results {
		status := "ok"
		if r.Err != nil {
			status = "failed"
		}
		fmt.Fprintf(w, "[hook:%s] %s (%s) — %dms\n",
			r.Hook.Stage,
			r.Hook.Command,
			status,
			r.Elapsed.Milliseconds(),
		)
		if r.Output != "" {
			fmt.Fprintf(w, "  output: %s\n", r.Output)
		}
		if r.Err != nil {
			fmt.Fprintf(w, "  error:  %v\n", r.Err)
		}
	}
}
