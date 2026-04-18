// Package cmd wires together the vaultpull pipeline.
//
// It provides:
//   - ParseFlags: parses CLI arguments into an Options struct.
//   - Run: orchestrates config loading, Vault secret fetching,
//     and writing secrets to a local .env file.
//
// Typical usage from main:
//
//	opts, err := cmd.ParseFlags(os.Args[1:])
//	if err != nil { ... }
//	err = cmd.Run(opts.Output, opts.Overwrite)
package cmd
