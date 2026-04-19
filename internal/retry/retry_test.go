package retry_test

import (
	"errors"
	"testing"
	"time"

	"github.com/example/vaultpull/internal/retry"
)

func TestDo_SuccessFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(retry.DefaultOptions(), func(attempt int) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnError(t *testing.T) {
	calls := 0
	opts := retry.Options{MaxAttempts: 3, Delay: 0, Multiplier: 1}
	err := retry.Do(opts, func(attempt int) error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retry, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	opts := retry.Options{MaxAttempts: 3, Delay: 0, Multiplier: 1}
	calls := 0
	err := retry.Do(opts, func(attempt int) error {
		calls++
		return errors.New("always fails")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, retry.ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ZeroAttemptsRunsOnce(t *testing.T) {
	calls := 0
	opts := retry.Options{MaxAttempts: 0, Delay: 0}
	retry.Do(opts, func(attempt int) error { //nolint
		calls++
		return nil
	})
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_BackoffIncreases(t *testing.T) {
	times := []time.Time{}
	opts := retry.Options{MaxAttempts: 3, Delay: 10 * time.Millisecond, Multiplier: 2.0}
	retry.Do(opts, func(attempt int) error { //nolint
		times = append(times, time.Now())
		return errors.New("fail")
	})
	if len(times) != 3 {
		t.Fatalf("expected 3 timestamps")
	}
	gap1 := times[1].Sub(times[0])
	gap2 := times[2].Sub(times[1])
	if gap2 < gap1 {
		t.Errorf("expected increasing backoff: gap1=%v gap2=%v", gap1, gap2)
	}
}
