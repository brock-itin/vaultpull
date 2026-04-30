// Package circuit implements a simple circuit breaker for protecting
// Vault API calls from cascading failures during outages or rate limiting.
package circuit

import (
	"errors"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing, requests blocked
	StateHalfOpen              // probing for recovery
)

// ErrOpen is returned when the circuit breaker is open.
var ErrOpen = errors.New("circuit breaker is open")

// Options configures the circuit breaker.
type Options struct {
	// MaxFailures is the number of consecutive failures before opening.
	MaxFailures int
	// Timeout is how long to wait before moving to half-open.
	Timeout time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxFailures: 5,
		Timeout:     30 * time.Second,
	}
}

// Breaker is a circuit breaker.
type Breaker struct {
	opts     Options
	mu       sync.Mutex
	state    State
	failures int
	openedAt time.Time
}

// New creates a new Breaker with the given options.
func New(opts Options) *Breaker {
	if opts.MaxFailures <= 0 {
		opts.MaxFailures = DefaultOptions().MaxFailures
	}
	if opts.Timeout <= 0 {
		opts.Timeout = DefaultOptions().Timeout
	}
	return &Breaker{opts: opts}
}

// Do executes fn if the circuit is closed or half-open.
// It records success or failure and transitions state accordingly.
func (b *Breaker) Do(fn func() error) error {
	b.mu.Lock()
	if b.state == StateOpen {
		if time.Since(b.openedAt) >= b.opts.Timeout {
			b.state = StateHalfOpen
		} else {
			b.mu.Unlock()
			return ErrOpen
		}
	}
	b.mu.Unlock()

	err := fn()

	b.mu.Lock()
	defer b.mu.Unlock()

	if err != nil {
		b.failures++
		if b.failures >= b.opts.MaxFailures || b.state == StateHalfOpen {
			b.state = StateOpen
			b.openedAt = time.Now()
		}
		return err
	}

	b.failures = 0
	b.state = StateClosed
	return nil
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Reset forces the breaker back to closed state.
func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = StateClosed
	b.failures = 0
}
