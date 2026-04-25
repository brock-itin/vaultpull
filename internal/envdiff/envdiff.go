// Package envdiff compares two .env files and reports what changed between them.
package envdiff

import (
	"fmt"
	"io"
	"sort"
)

// ChangeType describes the kind of change for a key.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
	Same    ChangeType = "same"
)

// Entry represents a single key's change status between two env maps.
type Entry struct {
	Key    string
	Old    string
	New    string
	Change ChangeType
}

// Result holds the full comparison output.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if any entry is not Same.
func (r Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Change != Same {
			return true
		}
	}
	return false
}

// Summary returns counts of each change type.
func (r Result) Summary() map[ChangeType]int {
	m := map[ChangeType]int{Added: 0, Removed: 0, Changed: 0, Same: 0}
	for _, e := range r.Entries {
		m[e.Change]++
	}
	return m
}

// Compare computes the diff between oldMap and newMap.
func Compare(oldMap, newMap map[string]string) Result {
	seen := make(map[string]bool)
	var entries []Entry

	for k, oldVal := range oldMap {
		seen[k] = true
		if newVal, ok := newMap[k]; ok {
			if oldVal == newVal {
				entries = append(entries, Entry{Key: k, Old: oldVal, New: newVal, Change: Same})
			} else {
				entries = append(entries, Entry{Key: k, Old: oldVal, New: newVal, Change: Changed})
			}
		} else {
			entries = append(entries, Entry{Key: k, Old: oldVal, New: "", Change: Removed})
		}
	}

	for k, newVal := range newMap {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Old: "", New: newVal, Change: Added})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Result{Entries: entries}
}

// Print writes a human-readable diff to w, masking values.
func Print(w io.Writer, r Result, maskValues bool) {
	if !r.HasChanges() {
		fmt.Fprintln(w, "No changes detected.")
		return
	}
	for _, e := range r.Entries {
		switch e.Change {
		case Added:
			fmt.Fprintf(w, "+ %s\n", e.Key)
		case Removed:
			fmt.Fprintf(w, "- %s\n", e.Key)
		case Changed:
			if maskValues {
				fmt.Fprintf(w, "~ %s (value changed)\n", e.Key)
			} else {
				fmt.Fprintf(w, "~ %s: %q -> %q\n", e.Key, e.Old, e.New)
			}
		}
	}
}
