// Package output provides utilities for rendering structured sync results
// to the terminal, including status summaries, key counts, and error reporting.
package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/vaultpull/internal/diff"
	"github.com/your-org/vaultpull/internal/rotate"
	"github.com/your-org/vaultpull/internal/validate"
)

// Options controls how output is rendered.
type Options struct {
	// Writer is the destination for output. Defaults to os.Stdout.
	Writer io.Writer
	// Quiet suppresses all non-error output.
	Quiet bool
	// NoColor disables ANSI color codes.
	NoColor bool
}

// DefaultOptions returns sensible output defaults.
func DefaultOptions() Options {
	return Options{
		Writer:  os.Stdout,
		Quiet:   false,
		NoColor: false,
	}
}

// SyncResult holds the aggregated result of a vault pull operation.
type SyncResult struct {
	Path        string
	OutputFile  string
	DiffResult  *diff.Result
	Rotation    *rotate.Report
	Validation  *validate.Report
	ElapsedMS   int64
}

// PrintSyncResult writes a human-readable summary of a sync operation.
func PrintSyncResult(r SyncResult, opts Options) {
	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}
	if opts.Quiet {
		return
	}

	green := colorFn(opts.NoColor, "\033[32m")
	yellow := colorFn(opts.NoColor, "\033[33m")
	red := colorFn(opts.NoColor, "\033[31m")
	reset := colorFn(opts.NoColor, "\033[0m")

	fmt.Fprintf(opts.Writer, "\n%s✔ Synced%s %s → %s\n",
		green(""), reset(""), r.Path, r.OutputFile)

	if r.DiffResult != nil {
		s := diff.Summary(*r.DiffResult)
		if diff.HasChanges(*r.DiffResult) {
			fmt.Fprintf(opts.Writer, "  %schanges:%s %s\n", yellow(""), reset(""), s)
		} else {
			fmt.Fprintf(opts.Writer, "  no changes detected\n")
		}
	}

	if r.Rotation != nil && rotate.HasStale(*r.Rotation) {
		stale := rotate.StaleKeys(*r.Rotation)
		fmt.Fprintf(opts.Writer, "  %s⚠ stale keys:%s %s\n",
			yellow(""), reset(""), strings.Join(stale, ", "))
	}

	if r.Validation != nil && validate.HasIssues(*r.Validation) {
		fmt.Fprintf(opts.Writer, "  %s✖ validation issues:%s\n", red(""), reset(""))
		for _, issue := range r.Validation.Issues {
			fmt.Fprintf(opts.Writer, "    - %s\n", issue)
		}
	}

	if r.ElapsedMS > 0 {
		fmt.Fprintf(opts.Writer, "  completed in %dms\n", r.ElapsedMS)
	}
}

// PrintError writes a formatted error message to the writer (or stderr).
func PrintError(err error, opts Options) {
	w := opts.Writer
	if w == nil {
		w = os.Stderr
	}
	red := colorFn(opts.NoColor, "\033[31m")
	reset := colorFn(opts.NoColor, "\033[0m")
	fmt.Fprintf(w, "%s✖ error:%s %v\n", red(""), reset(""), err)
}

// colorFn returns a function that wraps a string with the given ANSI code,
// or a no-op if noColor is true.
func colorFn(noColor bool, code string) func(string) string {
	if noColor {
		return func(s string) string { return s }
	}
	return func(s string) string {
		if s == "" {
			return code
		}
		return code + s + "\033[0m"
	}
}
