package circuit_test

import (
	"errors"
	"testing"
	"time"

	"github.com/vaultpull/internal/circuit"
)

var errFake = errors.New("fake error")

func TestDo_SuccessKeepsClosed(t *testing.T) {
	b := circuit.New(circuit.Options{MaxFailures: 3, Timeout: time.Second})
	if err := b.Do(func() error { return nil }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.State() != circuit.StateClosed {
		t.Fatalf("expected closed, got %v", b.State())
	}
}

func TestDo_OpensAfterMaxFailures(t *testing.T) {
	b := circuit.New(circuit.Options{MaxFailures: 3, Timeout: time.Second})
	for i := 0; i < 3; i++ {
		_ = b.Do(func() error { return errFake })
	}
	if b.State() != circuit.StateOpen {
		t.Fatalf("expected open, got %v", b.State())
	}
}

func TestDo_BlocksWhenOpen(t *testing.T) {
	b := circuit.New(circuit.Options{MaxFailures: 1, Timeout: time.Hour})
	_ = b.Do(func() error { return errFake })

	called := false
	err := b.Do(func() error { called = true; return nil })
	if !errors.Is(err, circuit.ErrOpen) {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
	if called {
		t.Fatal("fn should not have been called")
	}
}

func TestDo_HalfOpenAfterTimeout(t *testing.T) {
	b := circuit.New(circuit.Options{MaxFailures: 1, Timeout: 10 * time.Millisecond})
	_ = b.Do(func() error { return errFake })

	time.Sleep(20 * time.Millisecond)

	// Should be allowed through in half-open state.
	err := b.Do(func() error { return nil })
	if err != nil {
		t.Fatalf("unexpected error in half-open: %v", err)
	}
	if b.State() != circuit.StateClosed {
		t.Fatalf("expected closed after success in half-open, got %v", b.State())
	}
}

func TestDo_HalfOpenFailureReopens(t *testing.T) {
	b := circuit.New(circuit.Options{MaxFailures: 1, Timeout: 10 * time.Millisecond})
	_ = b.Do(func() error { return errFake })

	time.Sleep(20 * time.Millisecond)

	_ = b.Do(func() error { return errFake })
	if b.State() != circuit.StateOpen {
		t.Fatalf("expected open after half-open failure, got %v", b.State())
	}
}

func TestReset_ForcesClosed(t *testing.T) {
	b := circuit.New(circuit.Options{MaxFailures: 1, Timeout: time.Hour})
	_ = b.Do(func() error { return errFake })
	b.Reset()
	if b.State() != circuit.StateClosed {
		t.Fatalf("expected closed after reset, got %v", b.State())
	}
}

func TestDefaultOptions_Applied(t *testing.T) {
	b := circuit.New(circuit.Options{})
	for i := 0; i < 5; i++ {
		_ = b.Do(func() error { return errFake })
	}
	if b.State() != circuit.StateOpen {
		t.Fatalf("expected open with default max failures, got %v", b.State())
	}
}
