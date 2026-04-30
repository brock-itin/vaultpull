package batch_test

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/batch"
)

func makeSuccessFetch(data map[string]map[string]string) batch.FetchFunc {
	return func(_ context.Context, path string) (map[string]string, error) {
		if secrets, ok := data[path]; ok {
			return secrets, nil
		}
		return nil, fmt.Errorf("path not found: %s", path)
	}
}

func TestRun_AllSuccess(t *testing.T) {
	data := map[string]map[string]string{
		"secret/a": {"KEY_A": "val_a"},
		"secret/b": {"KEY_B": "val_b"},
	}
	results := batch.Run(context.Background(), []string{"secret/a", "secret/b"}, makeSuccessFetch(data), batch.DefaultOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for path %q: %v", r.Path, r.Err)
		}
	}
}

func TestRun_PartialError(t *testing.T) {
	data := map[string]map[string]string{
		"secret/a": {"KEY_A": "val_a"},
	}
	results := batch.Run(context.Background(), []string{"secret/a", "secret/missing"}, makeSuccessFetch(data), batch.DefaultOptions())
	if !batch.HasErrors(results) {
		t.Fatal("expected errors but HasErrors returned false")
	}
}

func TestRun_StopOnError(t *testing.T) {
	var calls int32
	fn := func(_ context.Context, path string) (map[string]string, error) {
		time.Sleep(5 * time.Millisecond)
		atomic.AddInt32(&calls, 1)
		if path == "secret/fail" {
			return nil, errors.New("forced failure")
		}
		return map[string]string{"K": "v"}, nil
	}
	opts := batch.Options{Concurrency: 1, StopOnError: true}
	paths := []string{"secret/fail", "secret/b", "secret/c", "secret/d"}
	results := batch.Run(context.Background(), paths, fn, opts)
	if !batch.HasErrors(results) {
		t.Fatal("expected at least one error")
	}
}

func TestRun_OrderPreserved(t *testing.T) {
	data := map[string]map[string]string{
		"secret/a": {"A": "1"},
		"secret/b": {"B": "2"},
		"secret/c": {"C": "3"},
	}
	paths := []string{"secret/a", "secret/b", "secret/c"}
	results := batch.Run(context.Background(), paths, makeSuccessFetch(data), batch.DefaultOptions())
	for i, p := range paths {
		if results[i].Path != p {
			t.Errorf("index %d: expected path %q, got %q", i, p, results[i].Path)
		}
	}
}

func TestMerge_CombinesSecrets(t *testing.T) {
	results := []batch.Result{
		{Path: "secret/a", Secrets: map[string]string{"A": "1", "SHARED": "from_a"}},
		{Path: "secret/b", Secrets: map[string]string{"B": "2", "SHARED": "from_b"}},
		{Path: "secret/c", Err: errors.New("skip me")},
	}
	merged := batch.Merge(results)
	if merged["A"] != "1" {
		t.Errorf("expected A=1, got %q", merged["A"])
	}
	if merged["B"] != "2" {
		t.Errorf("expected B=2, got %q", merged["B"])
	}
	if merged["SHARED"] != "from_b" {
		t.Errorf("expected SHARED=from_b (later wins), got %q", merged["SHARED"])
	}
	if len(merged) != 3 {
		t.Errorf("expected 3 keys, got %d", len(merged))
	}
}

func TestRun_EmptyPaths(t *testing.T) {
	fn := func(_ context.Context, _ string) (map[string]string, error) {
		return nil, nil
	}
	results := batch.Run(context.Background(), []string{}, fn, batch.DefaultOptions())
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
