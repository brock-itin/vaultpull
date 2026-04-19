package merge_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/merge"
)

func TestApply_StrategyLast(t *testing.T) {
	a := map[string]string{"FOO": "a", "BAR": "a"}
	b := map[string]string{"FOO": "b", "BAZ": "b"}

	result, err := merge.Apply([]map[string]string{a, b}, merge.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "b" {
		t.Errorf("expected FOO=b, got %s", result["FOO"])
	}
	if result["BAR"] != "a" {
		t.Errorf("expected BAR=a, got %s", result["BAR"])
	}
	if result["BAZ"] != "b" {
		t.Errorf("expected BAZ=b, got %s", result["BAZ"])
	}
}

func TestApply_StrategyFirst(t *testing.T) {
	a := map[string]string{"FOO": "first"}
	b := map[string]string{"FOO": "second"}

	result, err := merge.Apply([]map[string]string{a, b}, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "first" {
		t.Errorf("expected FOO=first, got %s", result["FOO"])
	}
}

func TestApply_StrategyError_NoConflict(t *testing.T) {
	a := map[string]string{"FOO": "a"}
	b := map[string]string{"BAR": "b"}

	result, err := merge.Apply([]map[string]string{a, b}, merge.StrategyError)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestApply_StrategyError_Conflict(t *testing.T) {
	a := map[string]string{"FOO": "a"}
	b := map[string]string{"FOO": "b"}

	_, err := merge.Apply([]map[string]string{a, b}, merge.StrategyError)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	ce, ok := err.(*merge.ConflictError)
	if !ok {
		t.Fatalf("expected *ConflictError, got %T", err)
	}
	if ce.Key != "FOO" {
		t.Errorf("expected conflict key FOO, got %s", ce.Key)
	}
}

func TestApply_Empty(t *testing.T) {
	result, err := merge.Apply(nil, merge.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}

func TestKeys_Union(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "3", "C": "4"}

	keys := merge.Keys([]map[string]string{a, b})
	if len(keys) != 3 {
		t.Errorf("expected 3 unique keys, got %d", len(keys))
	}
}
