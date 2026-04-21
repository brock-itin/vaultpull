package env

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ScanResult holds the result of scanning an env file for a specific key.
type ScanResult struct {
	Key      string
	Value    string
	Line     int
	Found    bool
	Comments []string // inline or preceding comments
}

// ScanOptions controls scanning behaviour.
type ScanOptions struct {
	CaseSensitive bool
}

// DefaultScanOptions returns sensible defaults.
func DefaultScanOptions() ScanOptions {
	return ScanOptions{
		CaseSensitive: true,
	}
}

// ScanFile opens the file at path and scans for the given key.
func ScanFile(path, key string, opts ScanOptions) (ScanResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return ScanResult{Key: key}, err
	}
	defer f.Close()
	return Scan(f, key, opts)
}

// Scan reads from r and searches for the given key, returning its position and
// value. Comments immediately preceding the key line are collected.
func Scan(r io.Reader, key string, opts ScanOptions) (ScanResult, error) {
	result := ScanResult{Key: key}
	var pending []string

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			pending = nil
			continue
		}

		if strings.HasPrefix(line, "#") {
			pending = append(pending, line)
			continue
		}

		eqIdx := strings.IndexByte(line, '=')
		if eqIdx < 0 {
			pending = nil
			continue
		}

		k := line[:eqIdx]
		v := strings.Trim(line[eqIdx+1:], "\"'")

		cmpKey := k
		cmpTarget := key
		if !opts.CaseSensitive {
			cmpKey = strings.ToUpper(k)
			cmpTarget = strings.ToUpper(key)
		}

		if cmpKey == cmpTarget {
			result.Found = true
			result.Value = v
			result.Line = lineNum
			result.Comments = pending
			return result, nil
		}

		pending = nil
	}

	return result, scanner.Err()
}
