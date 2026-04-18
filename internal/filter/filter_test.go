package filter_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/filter"
)

var base = map[string]string{
	"APP_HOST":  "localhost",
	"APP_PORT":  "8080",
	"DB_HOST":   "db.local",
	"DB_PASS":   "secret",
	"LOG_LEVEL": "info",
}

func TestApply_NoRule(t *testing.T) {
	out := filter.Apply(base, filter.Rule{})
	if len(out) != len(base) {
		t.Fatalf("expected %d keys, got %d", len(base), len(out))
	}
}

func TestApply_Prefix(t *testing.T) {
	out := filter.Apply(base, filter.Rule{Prefix: "APP_"})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
}

func TestApply_Only(t *testing.T) {
	out := filter.Apply(base, filter.Rule{Only: []string{"DB_HOST", "DB_PASS"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("APP_HOST should have been filtered out")
	}
}

func TestApply_Exclude(t *testing.T) {
	out := filter.Apply(base, filter.Rule{Exclude: []string{"DB_PASS", "LOG_LEVEL"}})
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should have been excluded")
	}
}

func TestApply_PrefixAndExclude(t *testing.T) {
	out := filter.Apply(base, filter.Rule{Prefix: "DB_", Exclude: []string{"DB_PASS"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestApply_OnlyOverridesPrefix(t *testing.T) {
	// Only takes precedence alongside prefix — both must match
	out := filter.Apply(base, filter.Rule{Prefix: "APP_", Only: []string{"DB_HOST"}})
	if len(out) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(out))
	}
}
