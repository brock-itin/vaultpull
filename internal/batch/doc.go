// Package batch provides concurrent secret fetching from multiple Vault paths.
//
// Use batch.Run to fan out fetch operations across a list of paths, controlling
// parallelism via Options.Concurrency. Results are returned in input order,
// making it safe to zip them back against the original path list.
//
// Example:
//
//	opts := batch.DefaultOptions()
//	opts.Concurrency = 8
//
//	results := batch.Run(ctx, paths, vaultClient.GetSecrets, opts)
//	if batch.HasErrors(results) {
//		// handle partial failures
//	}
//	merged := batch.Merge(results)
package batch
