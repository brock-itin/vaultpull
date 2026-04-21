package watch_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/watch"
)

// stubFetcher implements SecretFetcher for tests.
type stubFetcher struct {
	secrets map[string]string
	err     error
	calls   int
}

func (s *stubFetcher) GetSecrets(_ string) (map[string]string, error) {
	s.calls++
	return s.secrets, s.err
}

// stubWriter implements EnvWriter for tests.
type stubWriter struct {
	written map[string]string
	err     error
	calls   int
}

func (s *stubWriter) Write(_ string, secrets map[string]string) error {
	s.calls++
	s.written = secrets
	return s.err
}

func TestRun_InitialTickAndCancel(t *testing.T) {
	fetcher := &stubFetcher{secrets: map[string]string{"KEY": "val"}}
	writer := &stubWriter{}
	var buf bytes.Buffer

	opts := watch.Options{
		Interval:   50 * time.Millisecond,
		VaultPath:  "secret/app",
		OutputFile: ".env",
		Out:        &buf,
	}
	w := watch.New(opts, fetcher, writer)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if fetcher.calls < 1 {
		t.Errorf("expected at least 1 fetch call, got %d", fetcher.calls)
	}
	if writer.calls < 1 {
		t.Errorf("expected at least 1 write call, got %d", writer.calls)
	}
	if !strings.Contains(buf.String(), "synced") {
		t.Errorf("expected output to contain 'synced', got: %s", buf.String())
	}
}

func TestRun_FetchError_ContinuesLoop(t *testing.T) {
	fetcher := &stubFetcher{err: errors.New("vault unavailable")}
	writer := &stubWriter{}
	var buf bytes.Buffer

	opts := watch.Options{
		Interval:   40 * time.Millisecond,
		VaultPath:  "secret/app",
		OutputFile: ".env",
		Out:        &buf,
	}
	w := watch.New(opts, fetcher, writer)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Millisecond)
	defer cancel()

	w.Run(ctx) //nolint:errcheck

	if !strings.Contains(buf.String(), "fetch error") {
		t.Errorf("expected fetch error in output, got: %s", buf.String())
	}
	if writer.calls != 0 {
		t.Errorf("expected no write calls on fetch error, got %d", writer.calls)
	}
}

func TestNew_DefaultInterval(t *testing.T) {
	w := watch.New(watch.Options{Interval: 0}, &stubFetcher{}, &stubWriter{})
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}
