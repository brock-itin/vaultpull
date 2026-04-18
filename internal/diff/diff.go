// Package diff provides utilities for comparing secret maps
// against existing .env entries to detect additions, updates, and deletions.
package diff

// ChangeType represents the kind of change for a secret key.
type ChangeType string

const (
	Added   ChangeType = "added"
	Updated ChangeType = "updated"
	Deleted ChangeType = "deleted"
	Unchanged ChangeType = "unchanged"
)

// Change describes a single key-level diff result.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Compare returns the list of changes between existing and incoming secret maps.
// Keys present in incoming but not existing are Added.
// Keys present in both but with different values are Updated.
// Keys present in existing but not incoming are Deleted.
func Compare(existing, incoming map[string]string) []Change {
	var changes []Change

	for k, newVal := range incoming {
		oldVal, ok := existing[k]
		if !ok {
			changes = append(changes, Change{Key: k, Type: Added, NewValue: newVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Type: Updated, OldValue: oldVal, NewValue: newVal})
		} else {
			changes = append(changes, Change{Key: k, Type: Unchanged, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, oldVal := range existing {
		if _, ok := incoming[k]; !ok {
			changes = append(changes, Change{Key: k, Type: Deleted, OldValue: oldVal})
		}
	}

	return changes
}

// Summary returns counts of each change type from a list of changes.
func Summary(changes []Change) map[ChangeType]int {
	summary := map[ChangeType]int{}
	for _, c := range changes {
		summary[c.Type]++
	}
	return summary
}

// HasChanges returns true if any change in the list is not Unchanged.
func HasChanges(changes []Change) bool {
	for _, c := range changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}
