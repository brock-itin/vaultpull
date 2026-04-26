// Package pager provides pagination support for listing secrets from
// Vault paths that return large result sets. It wraps a fetch function
// and iterates through pages until all entries have been collected.
package pager

import (
	"context"
	"errors"
	"fmt"
)

// DefaultOptions returns sensible defaults for pagination.
func DefaultOptions() Options {
	return Options{
		PageSize: 100,
		MaxPages: 50,
	}
}

// Options controls pagination behaviour.
type Options struct {
	// PageSize is the number of items requested per page.
	PageSize int

	// MaxPages is the maximum number of pages to fetch before stopping,
	// preventing runaway iteration on unexpectedly large result sets.
	MaxPages int
}

// FetchFunc is a function that retrieves a single page of keys.
// offset is the zero-based starting index; limit is the page size.
// It must return the keys for the page and the total number of items
// available (used to determine whether more pages exist).
type FetchFunc func(ctx context.Context, offset, limit int) (keys []string, total int, err error)

// Result holds the aggregated output of a paginated fetch.
type Result struct {
	// Keys contains all collected keys across all pages.
	Keys []string

	// Pages is the number of pages fetched.
	Pages int

	// Truncated is true when MaxPages was reached before all items
	// were collected.
	Truncated bool
}

// Collect calls fetch repeatedly, advancing the offset by PageSize each
// time, until all items have been retrieved or MaxPages is reached.
func Collect(ctx context.Context, fetch FetchFunc, opts Options) (Result, error) {
	if opts.PageSize <= 0 {
		return Result{}, errors.New("pager: PageSize must be greater than zero")
	}
	if opts.MaxPages <= 0 {
		return Result{}, errors.New("pager: MaxPages must be greater than zero")
	}

	var result Result
	offset := 0

	for {
		if result.Pages >= opts.MaxPages {
			result.Truncated = true
			break
		}

		keys, total, err := fetch(ctx, offset, opts.PageSize)
		if err != nil {
			return Result{}, fmt.Errorf("pager: fetch at offset %d: %w", offset, err)
		}

		result.Keys = append(result.Keys, keys...)
		result.Pages++
		offset += len(keys)

		// Stop when we have received everything or the page was empty.
		if len(keys) == 0 || offset >= total {
			break
		}
	}

	return result, nil
}
