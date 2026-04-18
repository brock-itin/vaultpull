package audit

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileLogger wraps Logger with an underlying file handle.
type FileLogger struct {
	*Logger
	f *os.File
}

// NewFileLogger opens (or creates) a log file at path and returns a FileLogger.
func NewFileLogger(path string) (*FileLogger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("audit: mkdir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open: %w", err)
	}
	return &FileLogger{Logger: New(f), f: f}, nil
}

// Close closes the underlying log file.
func (fl *FileLogger) Close() error {
	if fl.f != nil {
		return fl.f.Close()
	}
	return nil
}
