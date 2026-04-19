package template

import (
	"fmt"
	"os"
	"path/filepath"
)

// RenderFile reads a template from srcPath, renders it with secrets, and
// writes the result to dstPath, creating parent directories as needed.
func RenderFile(srcPath, dstPath string, secrets map[string]string, opts *Options) error {
	raw, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("read template file: %w", err)
	}

	output, err := Render(string(raw), secrets, opts)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return fmt.Errorf("create output dirs: %w", err)
	}

	if err := os.WriteFile(dstPath, []byte(output), 0o600); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}
