// Package hook provides pre- and post-sync lifecycle hooks for vaultpull.
// Hooks are shell commands executed before or after secrets are written.
package hook

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Stage represents when a hook should run.
type Stage string

const (
	StagePre  Stage = "pre"
	StagePost Stage = "post"
)

// Hook defines a single lifecycle hook.
type Hook struct {
	Stage   Stage
	Command string
	Timeout time.Duration
}

// Result holds the outcome of running a hook.
type Result struct {
	Hook    Hook
	Output  string
	Err     error
	Elapsed time.Duration
}

// Runner executes lifecycle hooks.
type Runner struct {
	hooks []Hook
}

// New creates a Runner with the given hooks.
func New(hooks []Hook) *Runner {
	return &Runner{hooks: hooks}
}

// Run executes all hooks matching the given stage.
// It returns a slice of results, one per matching hook.
func (r *Runner) Run(ctx context.Context, stage Stage) []Result {
	var results []Result
	for _, h := range r.hooks {
		if h.Stage != stage {
			continue
		}
		results = append(results, r.exec(ctx, h))
	}
	return results
}

func (r *Runner) exec(ctx context.Context, h Hook) Result {
	timeout := h.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	parts := strings.Fields(h.Command)
	if len(parts) == 0 {
		return Result{Hook: h, Err: fmt.Errorf("empty command"), Elapsed: 0}
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	out, err := cmd.CombinedOutput()
	return Result{
		Hook:    h,
		Output:  strings.TrimSpace(string(out)),
		Err:     err,
		Elapsed: time.Since(start),
	}
}

// HasFailures returns true if any result contains an error.
func HasFailures(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}
