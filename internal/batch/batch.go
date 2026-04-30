// Package batch provides utilities for processing secrets in batches,
// allowing concurrent fetching from multiple Vault paths with configurable
// concurrency limits and error handling strategies.
package batch

import (
	"context"
	"fmt"
	"sync"
)

// Result holds the outcome of fetching secrets from a single path.
type Result struct {
	Path    string
	Secrets map[string]string
	Err     error
}

// FetchFunc is a function that retrieves secrets from a given path.
type FetchFunc func(ctx context.Context, path string) (map[string]string, error)

// Options configures batch processing behaviour.
type Options struct {
	// Concurrency is the maximum number of parallel fetches. Defaults to 4.
	Concurrency int
	// StopOnError causes the batch to abort remaining work on first error.
	StopOnError bool
}

// DefaultOptions returns sensible defaults for batch processing.
func DefaultOptions() Options {
	return Options{
		Concurrency: 4,
		StopOnError: false,
	}
}

// Run fetches secrets from all provided paths concurrently, respecting the
// configured concurrency limit. Results are returned in the same order as
// the input paths.
func Run(ctx context.Context, paths []string, fn FetchFunc, opts Options) []Result {
	if opts.Concurrency <= 0 {
		opts.Concurrency = DefaultOptions().Concurrency
	}

	results := make([]Result, len(paths))
	sem := make(chan struct{}, opts.Concurrency)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for i, path := range paths {
		wg.Add(1)
		go func(idx int, p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if ctx.Err() != nil {
				results[idx] = Result{Path: p, Err: ctx.Err()}
				return
			}

			secrets, err := fn(ctx, p)
			results[idx] = Result{Path: p, Secrets: secrets, Err: err}

			if err != nil && opts.StopOnError {
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("batch: path %q: %w", p, err)
					cancel()
				}
				mu.Unlock()
			}
		}(i, path)
	}

	wg.Wait()
	return results
}

// HasErrors reports whether any result in the slice contains an error.
func HasErrors(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}

// Merge combines all successful results into a single map. Later paths
// overwrite earlier ones on key conflicts.
func Merge(results []Result) map[string]string {
	out := make(map[string]string)
	for _, r := range results {
		if r.Err != nil {
			continue
		}
		for k, v := range r.Secrets {
			out[k] = v
		}
	}
	return out
}
