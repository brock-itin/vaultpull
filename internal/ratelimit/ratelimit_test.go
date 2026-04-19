package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/ratelimit"
)

func TestNew_InvalidRate(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Options{RequestsPerSecond: 0})
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
}

func TestNew_Valid(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Options{RequestsPerSecond: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	// 1 rps — first call consumes the token, second must wait
	l, err := ratelimit.New(ratelimit.Options{RequestsPerSecond: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx := context.Background()
	// consume the initial token
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("first wait failed: %v", err)
	}
	// second call should block; cancel quickly
	ctx2, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err = l.Wait(ctx2)
	if err == nil {
		t.Fatal("expected context deadline error")
	}
}

func TestWait_HighRate_DoesNotBlock(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Options{RequestsPerSecond: 1000})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 10; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("wait %d failed: %v", i, err)
		}
	}
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Fatalf("10 requests at 1000 rps took too long: %v", elapsed)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := ratelimit.DefaultOptions()
	if opts.RequestsPerSecond <= 0 {
		t.Fatalf("expected positive default rate, got %d", opts.RequestsPerSecond)
	}
}
