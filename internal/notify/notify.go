// Package notify provides a simple notification dispatcher for vaultpull
// sync events. It supports multiple channels (stdout, webhook, and file)
// so operators can react to sync completions, errors, and drift detections
// without coupling the core sync logic to any specific alerting backend.
package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Level represents the severity of a notification event.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Event holds the data for a single notification.
type Event struct {
	Level   Level             `json:"level"`
	Message string            `json:"message"`
	Path    string            `json:"path,omitempty"`
	Meta    map[string]string `json:"meta,omitempty"`
	Time    time.Time         `json:"time"`
}

// Channel is the interface implemented by every notification backend.
type Channel interface {
	Send(ctx context.Context, e Event) error
}

// Dispatcher fans an event out to all registered channels, collecting
// any errors that occur along the way.
type Dispatcher struct {
	channels []Channel
}

// New returns a Dispatcher wired to the provided channels.
func New(channels ...Channel) *Dispatcher {
	return &Dispatcher{channels: channels}
}

// Send delivers the event to every channel. It always stamps the event
// with the current UTC time if the caller left Time at the zero value.
func (d *Dispatcher) Send(ctx context.Context, e Event) []error {
	if e.Time.IsZero() {
		e.Time = time.Now().UTC()
	}
	var errs []error
	for _, ch := range d.channels {
		if err := ch.Send(ctx, e); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// ─── Stdout channel ───────────────────────────────────────────────────────────

// StdoutChannel writes human-readable events to an io.Writer (usually os.Stdout).
type StdoutChannel struct {
	Out io.Writer
}

// NewStdoutChannel returns a StdoutChannel that writes to os.Stdout.
func NewStdoutChannel() *StdoutChannel {
	return &StdoutChannel{Out: os.Stdout}
}

func (c *StdoutChannel) Send(_ context.Context, e Event) error {
	_, err := fmt.Fprintf(c.Out, "[%s] %s %s\n",
		e.Level, e.Time.Format(time.RFC3339), e.Message)
	return err
}

// ─── Webhook channel ─────────────────────────────────────────────────────────

// WebhookChannel POSTs JSON-encoded events to a remote URL.
type WebhookChannel struct {
	URL    string
	Client *http.Client
}

// NewWebhookChannel returns a WebhookChannel with a sensible default timeout.
func NewWebhookChannel(url string) *WebhookChannel {
	return &WebhookChannel{
		URL:    url,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *WebhookChannel) Send(ctx context.Context, e Event) error {
	body, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("notify: marshal event: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("notify: send webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: webhook returned status %d", resp.StatusCode)
	}
	return nil
}

// ─── File channel ─────────────────────────────────────────────────────────────

// FileChannel appends JSON-encoded events to a log file, one per line.
type FileChannel struct {
	Path string
}

// NewFileChannel returns a FileChannel that appends to the given path.
func NewFileChannel(path string) *FileChannel {
	return &FileChannel{Path: path}
}

func (c *FileChannel) Send(_ context.Context, e Event) error {
	f, err := os.OpenFile(c.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("notify: open log file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("notify: write log file: %w", err)
	}
	return nil
}
