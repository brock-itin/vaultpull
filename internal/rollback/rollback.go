// Package rollback provides functionality to restore a .env file
// from a previously saved snapshot, enabling safe recovery after
// a failed or undesired sync operation.
package rollback

import (
	"fmt"
	"os"
	"time"

	"github.com/vaultpull/internal/snapshot"
)

// Result describes the outcome of a rollback operation.
type Result struct {
	Path        string
	SnapshotAt  time.Time
	KeysRestored int
	Err         error
}

// Options configures rollback behaviour.
type Options struct {
	// SnapshotDir is the directory where snapshots are stored.
	SnapshotDir string
	// DryRun reports what would be restored without writing anything.
	DryRun bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		SnapshotDir: ".vaultpull/snapshots",
	}
}

// Execute loads the snapshot for the given env file path and writes it
// back to disk, effectively undoing the last sync.
func Execute(envPath string, opts Options) Result {
	snapshotPath := snapshotFile(opts.SnapshotDir, envPath)

	entries, err := snapshot.Load(snapshotPath)
	if err != nil {
		return Result{Path: envPath, Err: fmt.Errorf("load snapshot: %w", err)}
	}
	if len(entries) == 0 {
		return Result{Path: envPath, Err: fmt.Errorf("no snapshot found for %s", envPath)}
	}

	data := snapshot.ToMap(entries)

	if opts.DryRun {
		return Result{
			Path:         envPath,
			KeysRestored: len(data),
		}
	}

	if err := writeEnvFile(envPath, data); err != nil {
		return Result{Path: envPath, Err: fmt.Errorf("write env file: %w", err)}
	}

	var snapshotAt time.Time
	if len(entries) > 0 {
		snapshotAt = entries[0].SavedAt
	}

	return Result{
		Path:         envPath,
		SnapshotAt:   snapshotAt,
		KeysRestored: len(data),
	}
}

func snapshotFile(dir, envPath string) string {
	safe := sanitizePath(envPath)
	return fmt.Sprintf("%s/%s.snap", dir, safe)
}

func sanitizePath(p string) string {
	out := make([]byte, len(p))
	for i := range p {
		if p[i] == '/' || p[i] == '.' {
			out[i] = '_'
		} else {
			out[i] = p[i]
		}
	}
	return string(out)
}

func writeEnvFile(path string, data map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for k, v := range data {
		if _, err := fmt.Fprintf(f, "%s=%s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}
