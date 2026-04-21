package rollback_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vaultpull/internal/rollback"
	"github.com/vaultpull/internal/snapshot"
)

func writeSnapshot(t *testing.T, dir, envPath string, data map[string]string) {
	t.Helper()
	entries := snapshot.FromMap(data, time.Now())
	safe := sanitize(envPath)
	path := filepath.Join(dir, safe+".snap")
	if err := snapshot.Save(path, entries); err != nil {
		t.Fatalf("setup: save snapshot: %v", err)
	}
}

func sanitize(p string) string {
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

func TestExecute_Success(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")
	snapshotDir := filepath.Join(tmpDir, "snapshots")

	data := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc123"}
	writeSnapshot(t, snapshotDir, envFile, data)

	opts := rollback.Options{SnapshotDir: snapshotDir}
	res := rollback.Execute(envFile, opts)

	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.KeysRestored != 2 {
		t.Errorf("expected 2 keys restored, got %d", res.KeysRestored)
	}
	if _, err := os.Stat(envFile); err != nil {
		t.Errorf("env file not created: %v", err)
	}
}

func TestExecute_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")
	snapshotDir := filepath.Join(tmpDir, "snapshots")

	data := map[string]string{"FOO": "bar"}
	writeSnapshot(t, snapshotDir, envFile, data)

	opts := rollback.Options{SnapshotDir: snapshotDir, DryRun: true}
	res := rollback.Execute(envFile, opts)

	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.KeysRestored != 1 {
		t.Errorf("expected 1 key, got %d", res.KeysRestored)
	}
	if _, err := os.Stat(envFile); !os.IsNotExist(err) {
		t.Error("dry run should not write env file")
	}
}

func TestExecute_MissingSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	opts := rollback.Options{SnapshotDir: filepath.Join(tmpDir, "snapshots")}
	res := rollback.Execute(envFile, opts)

	if res.Err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := rollback.DefaultOptions()
	if opts.SnapshotDir == "" {
		t.Error("expected non-empty default SnapshotDir")
	}
	if opts.DryRun {
		t.Error("DryRun should default to false")
	}
}
