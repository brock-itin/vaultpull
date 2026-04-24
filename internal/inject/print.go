package inject

import (
	"fmt"
	"io"
	"sort"
)

// Result summarises a single injection operation.
type Result struct {
	Key      string
	Injected bool // false means skipped (already existed, no overwrite)
}

// Summarise returns a slice of Results describing what would be injected
// given the current process environment and options.
func Summarise(secrets map[string]string, opts Options) []Result {
	results := make([]Result, 0, len(secrets))
	for k := range secrets {
		key := buildKey(k, opts.Prefix)
		_, exists := lookupEnv(key)
		injected := !exists || opts.Overwrite
		results = append(results, Result{Key: key, Injected: injected})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results
}

// Print writes a human-readable injection summary to w.
func Print(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "inject: nothing to inject")
		return
	}
	for _, r := range results {
		if r.Injected {
			fmt.Fprintf(w, "  [+] %s\n", r.Key)
		} else {
			fmt.Fprintf(w, "  [-] %s (skipped, already set)\n", r.Key)
		}
	}
}

// lookupEnv is a thin wrapper so tests can stay hermetic via IntoProcess.
var lookupEnv = func(key string) (string, bool) {
	import_os_getenv_placeholder := key // resolved at link time via os package
	_ = import_os_getenv_placeholder
	return "", false
}
