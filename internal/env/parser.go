package env

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair parsed from a .env file.
type Entry struct {
	Key   string
	Value string
	Raw   string // original line, preserved for comments/blanks
}

// Parse reads an .env file from r and returns a slice of entries.
// Comment lines and blank lines are preserved as entries with empty Key.
func Parse(r io.Reader) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Blank lines or comments — preserve as-is
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			entries = append(entries, Entry{Raw: line})
			continue
		}

		idx := strings.IndexByte(trimmed, '=')
		if idx < 0 {
			return nil, fmt.Errorf("parse error on line %d: missing '=' in %q", lineNum, line)
		}

		key := strings.TrimSpace(trimmed[:idx])
		val := strings.TrimSpace(trimmed[idx+1:])

		// Strip surrounding quotes if present
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		if key == "" {
			return nil, fmt.Errorf("parse error on line %d: empty key", lineNum)
		}

		entries = append(entries, Entry{Key: key, Value: val, Raw: line})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return entries, nil
}

// ToMap converts a slice of entries into a key→value map.
// Comment/blank entries (empty Key) are skipped.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.Key != "" {
			m[e.Key] = e.Value
		}
	}
	return m
}
