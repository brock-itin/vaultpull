package diff

import (
	"fmt"
	"io"
	"sort"
)

// Print writes a human-readable diff summary to w.
// Secret values are masked to avoid leaking sensitive data.
func Print(w io.Writer, changes []Change) {
	if len(changes) == 0 {
		fmt.Fprintln(w, "No changes detected.")
		return
	}

	sorted := make([]Change, len(changes))
	copy(sorted, changes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, c := range sorted {
		switch c.Type {
		case Added:
			fmt.Fprintf(w, "  + %s\n", c.Key)
		case Updated:
			fmt.Fprintf(w, "  ~ %s (changed)\n", c.Key)
		case Deleted:
			fmt.Fprintf(w, "  - %s\n", c.Key)
		case Unchanged:
			fmt.Fprintf(w, "    %s (unchanged)\n", c.Key)
		}
	}

	s := Summary(changes)
	fmt.Fprintf(w, "\nSummary: +%d added, ~%d updated, -%d deleted, %d unchanged\n",
		s[Added], s[Updated], s[Deleted], s[Unchanged])
}
