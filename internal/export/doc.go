// Package export provides multi-format secret export for vaultpull.
//
// It supports writing secrets fetched from Vault into local files using
// several output formats:
//
//   - dotenv  — standard KEY="value" pairs (default)
//   - shell   — export KEY='value' lines suitable for sourcing
//   - docker  — plain KEY=value lines for use with --env-file
//   - json    — a JSON object of key/value pairs
//
// # Basic usage
//
//	opts := export.DefaultOptions()
//	opts.Format = export.FormatShell
//	opts.OutputPath = "/tmp/secrets.sh"
//	export.Export(secrets, opts)
//
// # Multi-target export
//
// Targets can be defined as strings and parsed with ParseTarget:
//
//	tgt, _ := export.ParseTarget("ci:/tmp/ci.env@docker+APP_")
//	export.ExportAll(secrets, []export.Target{tgt}, export.DefaultOptions())
package export
