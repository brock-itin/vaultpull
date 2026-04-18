// Package audit provides structured audit logging for vaultpull operations.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Path      string    `json:"path,omitempty"`
	Output    string    `json:"output,omitempty"`
	Keys      []string  `json:"keys,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes audit entries to a destination.
type Logger struct {
	w io.Writer
}

// New returns a Logger writing to w. Pass nil to use os.Stderr.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{w: w}
}

// Log writes a JSON-encoded audit entry.
func (l *Logger) Log(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	if err != nil {
		return fmt.Errorf("audit: write: %w", err)
	}
	return nil
}

// Success logs a successful sync event.
func (l *Logger) Success(path, output string, keys []string) error {
	return l.Log(Entry{
		Event:  "sync_success",
		Path:   path,
		Output: output,
		Keys:   keys,
	})
}

// Failure logs a failed sync event.
func (l *Logger) Failure(path string, err error) error {
	return l.Log(Entry{
		Event: "sync_failure",
		Path:  path,
		Error: err.Error(),
	})
}
