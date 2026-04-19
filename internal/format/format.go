// Package format provides output formatters for secret maps.
package format

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Type represents an output format.
type Type string

const (
	TypeEnv  Type = "env"
	TypeJSON Type = "json"
	TypeExport Type = "export"
)

// Write writes secrets in the given format to w.
func Write(w io.Writer, secrets map[string]string, t Type) error {
	switch t {
	case TypeJSON:
		return writeJSON(w, secrets)
	case TypeExport:
		return writeExport(w, secrets)
	case TypeEnv:
		return writeEnv(w, secrets)
	default:
		return fmt.Errorf("unknown format: %q", t)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func writeEnv(w io.Writer, secrets map[string]string) error {
	for _, k := range sortedKeys(secrets) {
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, secrets[k]); err != nil {
			return err
		}
	}
	return nil
}

func writeExport(w io.Writer, secrets map[string]string) error {
	for _, k := range sortedKeys(secrets) {
		v := strings.ReplaceAll(secrets[k], "'", "'\"'\"'")
		if _, err := fmt.Fprintf(w, "export %s='%s'\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, secrets map[string]string) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(secrets)
}
