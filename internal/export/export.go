// Package export provides functionality for exporting secrets to various
// output targets such as shell scripts, Docker env files, and JSON configs.
package export

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Format represents the output format for exported secrets.
type Format string

const (
	FormatShell  Format = "shell"
	FormatDocker Format = "docker"
	FormatJSON   Format = "json"
	FormatDotEnv Format = "dotenv"
)

// Options controls the behaviour of an export operation.
type Options struct {
	Format      Format
	OutputPath  string
	Prefix      string
	OmitEmpty   bool
	Permissions os.FileMode
}

// DefaultOptions returns sensible defaults for export operations.
func DefaultOptions() Options {
	return Options{
		Format:      FormatDotEnv,
		Permissions: 0600,
	}
}

// Export writes the provided secrets map to the configured target.
func Export(secrets map[string]string, opts Options) error {
	if opts.Permissions == 0 {
		opts.Permissions = 0600
	}

	if opts.OutputPath != "" {
		if err := os.MkdirAll(filepath.Dir(opts.OutputPath), 0755); err != nil {
			return fmt.Errorf("export: create parent dirs: %w", err)
		}
		f, err := os.OpenFile(opts.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, opts.Permissions)
		if err != nil {
			return fmt.Errorf("export: open file: %w", err)
		}
		defer f.Close()
		return write(f, secrets, opts)
	}

	return write(os.Stdout, secrets, opts)
}

func write(w io.Writer, secrets map[string]string, opts Options) error {
	keys := sortedKeys(secrets)

	switch opts.Format {
	case FormatShell:
		return writeShell(w, keys, secrets, opts)
	case FormatDocker:
		return writeDocker(w, keys, secrets, opts)
	case FormatJSON:
		return writeJSON(w, keys, secrets, opts)
	default:
		return writeDotEnv(w, keys, secrets, opts)
	}
}

func writeShell(w io.Writer, keys []string, secrets map[string]string, opts Options) error {
	for _, k := range keys {
		v := secrets[k]
		if opts.OmitEmpty && v == "" {
			continue
		}
		key := applyPrefix(k, opts.Prefix)
		v = strings.ReplaceAll(v, "'", "'\"'\"'")
		fmt.Fprintf(w, "export %s='%s'\n", key, v)
	}
	return nil
}

func writeDocker(w io.Writer, keys []string, secrets map[string]string, opts Options) error {
	for _, k := range keys {
		v := secrets[k]
		if opts.OmitEmpty && v == "" {
			continue
		}
		key := applyPrefix(k, opts.Prefix)
		fmt.Fprintf(w, "%s=%s\n", key, v)
	}
	return nil
}

func writeJSON(w io.Writer, keys []string, secrets map[string]string, opts Options) error {
	fmt.Fprintln(w, "{")
	for i, k := range keys {
		v := secrets[k]
		if opts.OmitEmpty && v == "" {
			continue
		}
		key := applyPrefix(k, opts.Prefix)
		v = strings.ReplaceAll(v, `"`, `\"`)
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(w, "  \"%s\": \"%s\"%s\n", key, v, comma)
	}
	fmt.Fprintln(w, "}")
	return nil
}

func writeDotEnv(w io.Writer, keys []string, secrets map[string]string, opts Options) error {
	for _, k := range keys {
		v := secrets[k]
		if opts.OmitEmpty && v == "" {
			continue
		}
		key := applyPrefix(k, opts.Prefix)
		fmt.Fprintf(w, "%s=%q\n", key, v)
	}
	return nil
}

func applyPrefix(key, prefix string) string {
	if prefix == "" {
		return key
	}
	return prefix + key
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
