// Package chain implements a composable pipeline for processing secret maps.
//
// A Pipeline is built by chaining named steps, each of which receives the
// current secret map and returns a (possibly modified) map or nil to skip.
// Steps are executed in order; any error halts the pipeline immediately and
// is returned wrapped with the failing step's name.
//
// Example usage:
//
//	result, err := chain.New().
//	    Add("filter",    filter.Apply).
//	    Add("transform", transform.Apply).
//	    Add("redact",    redact.Map).
//	    Run(secrets)
package chain
