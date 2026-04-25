package quota_test

import (
	"errors"
	"testing"
	"time"

	"github.com/vaultpull/internal/quota"
)

func TestRecord_UnderLimit(t *testing.T) {
	tr := quota.New(quota.Options{MaxFetches: 3, Window: time.Minute})
	for i := 0; i < 3; i++ {
		if err := tr.Record("secret/app"); err != nil {
			t.Fatalf("unexpected error on fetch %d: %v", i+1, err)
		}
	}
}

func TestRecord_ExceedsLimit(t *testing.T) {
	tr := quota.New(quota.Options{MaxFetches: 2, Window: time.Minute})
	_ = tr.Record("secret/app")
	_ = tr.Record("secret/app")

	err := tr.Record("secret/app")
	if err == nil {
		t.Fatal("expected quota exceeded error, got nil")
	}
	if !errors.Is(err, quota.ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got: %v", err)
	}
}

func TestRecord_IndependentPaths(t *testing.T) {
	tr := quota.New(quota.Options{MaxFetches: 1, Window: time.Minute})
	if err := tr.Record("secret/a"); err != nil {
		t.Fatalf("unexpected error for path a: %v", err)
	}
	if err := tr.Record("secret/b"); err != nil {
		t.Fatalf("unexpected error for path b: %v", err)
	}
}

func TestRecord_WindowExpiry(t *testing.T) {
	tr := quota.New(quota.Options{MaxFetches: 1, Window: 10 * time.Millisecond})
	_ = tr.Record("secret/app")

	time.Sleep(20 * time.Millisecond)

	// Window should have reset — this should succeed.
	if err := tr.Record("secret/app"); err != nil {
		t.Fatalf("expected window reset, got error: %v", err)
	}
}

func TestUsage_ReturnsCount(t *testing.T) {
	tr := quota.New(quota.Options{MaxFetches: 10, Window: time.Minute})
	_ = tr.Record("secret/db")
	_ = tr.Record("secret/db")

	if got := tr.Usage("secret/db"); got != 2 {
		t.Fatalf("expected usage 2, got %d", got)
	}
}

func TestUsage_UnknownPath(t *testing.T) {
	tr := quota.New(quota.Options{MaxFetches: 10, Window: time.Minute})
	if got := tr.Usage("secret/unknown"); got != 0 {
		t.Fatalf("expected 0 for unknown path, got %d", got)
	}
}

func TestReset_ClearsCounts(t *testing.T) {
	tr := quota.New(quota.Options{MaxFetches: 1, Window: time.Minute})
	_ = tr.Record("secret/app")
	tr.Reset()

	if got := tr.Usage("secret/app"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	// Should also allow fetching again after reset.
	if err := tr.Record("secret/app"); err != nil {
		t.Fatalf("unexpected error after reset: %v", err)
	}
}

func TestDefaultOptions_Sane(t *testing.T) {
	opts := quota.DefaultOptions()
	if opts.MaxFetches <= 0 {
		t.Errorf("expected positive MaxFetches, got %d", opts.MaxFetches)
	}
	if opts.Window <= 0 {
		t.Errorf("expected positive Window, got %v", opts.Window)
	}
}
