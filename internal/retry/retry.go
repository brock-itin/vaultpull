// Package retry provides simple retry logic for transient errors.
package retry

import (
	"errors"
	"time"
)

// Options configures retry behaviour.
type Options struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// DefaultOptions returns sensible retry defaults.
func DefaultOptions() Options {
	return Options{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// ErrMaxAttempts is returned when all attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Do calls fn up to opts.MaxAttempts times, backing off between attempts.
// fn should return a non-nil error to trigger a retry.
func Do(opts Options, fn func(attempt int) error) error {
	if opts.MaxAttempts <= 0 {
		opts.MaxAttempts = 1
	}
	delay := opts.Delay
	for i := 1; i <= opts.MaxAttempts; i++ {
		err := fn(i)
		if err == nil {
			return nil
		}
		if i == opts.MaxAttempts {
			return errors.Join(ErrMaxAttempts, err)
		}
		time.Sleep(delay)
		if opts.Multiplier > 0 {
			delay = time.Duration(float64(delay) * opts.Multiplier)
		}
	}
	return ErrMaxAttempts
}
