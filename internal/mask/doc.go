// Package mask provides helpers for redacting sensitive secret values
// before they appear in terminal output, audit logs, or diff displays.
//
// Use Value to mask a single string and Map to mask an entire secrets map.
// The Options type allows revealing a short prefix to aid debugging while
// keeping the bulk of the value hidden.
package mask
