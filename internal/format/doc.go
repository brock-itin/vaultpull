// Package format provides output formatters for writing secret maps
// in different representations: env file, shell export statements, and JSON.
//
// Supported formats:
//
//   - TypeEnv    — KEY=VALUE lines suitable for .env files
//   - TypeExport — export KEY='VALUE' lines for shell sourcing
//   - TypeJSON   — pretty-printed JSON object
package format
