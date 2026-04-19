// Package merge provides strategies for combining secret maps from multiple sources.
package merge

// Strategy defines how conflicts are resolved when merging secret maps.
type Strategy int

const (
	// StrategyFirst keeps the first value seen for a key.
	StrategyFirst Strategy = iota
	// StrategyLast overwrites with the last value seen for a key.
	StrategyLast
	// StrategyError returns an error on any conflicting key.
	StrategyError
)

// ConflictError is returned when StrategyError is used and a duplicate key is found.
type ConflictError struct {
	Key string
}

func (e *ConflictError) Error() string {
	return "merge conflict: duplicate key \"" + e.Key + "\""
}

// Apply merges a slice of secret maps according to the given strategy.
// Maps are applied in order; index 0 has lowest priority for StrategyLast.
func Apply(maps []map[string]string, strategy Strategy) (map[string]string, error) {
	result := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			_, exists := result[k]
			switch {
			case !exists:
				result[k] = v
			case strategy == StrategyFirst:
				// keep existing
			case strategy == StrategyLast:
				result[k] = v
			case strategy == StrategyError:
				return nil, &ConflictError{Key: k}
			}
		}
	}

	return result, nil
}

// Keys returns the union of all keys across the provided maps.
func Keys(maps []map[string]string) []string {
	seen := make(map[string]struct{})
	for _, m := range maps {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}
