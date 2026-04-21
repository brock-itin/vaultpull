// Package watch provides polling-based secret refresh for vaultpull.
// It periodically re-fetches secrets from Vault and writes updated .env files.
package watch

import (
	"context"
	"fmt"
	"io"
	"time"
)

// SecretFetcher retrieves secrets from a remote source.
type SecretFetcher interface {
	GetSecrets(path string) (map[string]string, error)
}

// EnvWriter writes a map of secrets to a local .env file.
type EnvWriter interface {
	Write(path string, secrets map[string]string) error
}

// Options configures the watch loop.
type Options struct {
	Interval  time.Duration
	VaultPath string
	OutputFile string
	Out       io.Writer
}

// DefaultOptions returns sensible defaults for the watch loop.
func DefaultOptions() Options {
	return Options{
		Interval: 30 * time.Second,
	}
}

// Watcher polls Vault for secret changes and rewrites the .env file.
type Watcher struct {
	opts    Options
	fetcher SecretFetcher
	writer  EnvWriter
}

// New creates a new Watcher with the given options and dependencies.
func New(opts Options, fetcher SecretFetcher, writer EnvWriter) *Watcher {
	if opts.Interval <= 0 {
		opts.Interval = DefaultOptions().Interval
	}
	if opts.Out == nil {
		opts.Out = io.Discard
	}
	return &Watcher{opts: opts, fetcher: fetcher, writer: writer}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	fmt.Fprintf(w.opts.Out, "[watch] starting — interval=%s path=%s\n",
		w.opts.Interval, w.opts.VaultPath)

	if err := w.tick(); err != nil {
		fmt.Fprintf(w.opts.Out, "[watch] initial fetch error: %v\n", err)
	}

	ticker := time.NewTicker(w.opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintf(w.opts.Out, "[watch] stopped\n")
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(); err != nil {
				fmt.Fprintf(w.opts.Out, "[watch] fetch error: %v\n", err)
			}
		}
	}
}

func (w *Watcher) tick() error {
	secrets, err := w.fetcher.GetSecrets(w.opts.VaultPath)
	if err != nil {
		return fmt.Errorf("fetch secrets: %w", err)
	}
	if err := w.writer.Write(w.opts.OutputFile, secrets); err != nil {
		return fmt.Errorf("write env: %w", err)
	}
	fmt.Fprintf(w.opts.Out, "[watch] synced %d secrets -> %s\n",
		len(secrets), w.opts.OutputFile)
	return nil
}
