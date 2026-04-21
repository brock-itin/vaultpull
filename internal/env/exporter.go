package env

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ExportOptions controls how secrets are exported to shell-sourceable files.
type ExportOptions struct {
	// Prefix is prepended to every key before exporting.
	Prefix string
	// Overwrite controls whether an existing file is replaced.
	Overwrite bool
	// Perm is the file permission used when creating the output file.
	Perm os.FileMode
}

// DefaultExportOptions returns sensible defaults for ExportOptions.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Overwrite: true,
		Perm:      0600,
	}
}

// Export writes secrets as "export KEY=VALUE" lines to the given path,
// creating parent directories as needed. If opts is nil, defaults are used.
func Export(path string, secrets map[string]string, opts *ExportOptions) error {
	if opts == nil {
		d := DefaultExportOptions()
		opts = &d
	}

	if !opts.Overwrite {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("export: file already exists: %s", path)
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("export: create dirs: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, opts.Perm)
	if err != nil {
		return fmt.Errorf("export: open file: %w", err)
	}
	defer f.Close()

	return writeExportLines(f, secrets, opts.Prefix)
}

func writeExportLines(w io.Writer, secrets map[string]string, prefix string) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		key := prefix + k
		val := strings.ReplaceAll(secrets[k], "'", "'\\''")
		if _, err := fmt.Fprintf(w, "export %s='%s'\n", key, val); err != nil {
			return fmt.Errorf("export: write key %s: %w", key, err)
		}
	}
	return nil
}
