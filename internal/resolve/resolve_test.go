package resolve

import (
	"testing"
)

func TestResolve_FlatKeys(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"DB_HOST": "localhost", "DB_PORT": "5432"},
	}
	opts := DefaultOptions()
	results, err := Resolve(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	m := ToMap(results)
	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", m["DB_HOST"])
	}
	if m["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", m["DB_PORT"])
	}
}

func TestResolve_WithPathPrefix(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app/db": {"HOST": "db.internal"},
	}
	opts := DefaultOptions()
	opts.StripPathPrefix = false

	results, err := Resolve(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "secret_app_db_HOST" {
		t.Errorf("expected key secret_app_db_HOST, got %q", results[0].Key)
	}
}

func TestResolve_SourcePopulated(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/svc": {"TOKEN": "abc123"},
	}
	results, err := Resolve(secrets, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Source != "secret/svc" {
		t.Errorf("expected source secret/svc, got %q", results[0].Source)
	}
}

func TestResolve_EmptyKeyReturnsError(t *testing.T) {
	secrets := map[string]map[string]string{
		"secret/app": {"": "value"},
	}
	_, err := Resolve(secrets, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestResolve_EmptySecrets(t *testing.T) {
	results, err := Resolve(map[string]map[string]string{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestToMap_Deduplicates(t *testing.T) {
	results := []Result{
		{Key: "FOO", Value: "first", Source: "a"},
		{Key: "FOO", Value: "second", Source: "b"},
	}
	m := ToMap(results)
	if len(m) != 1 {
		t.Errorf("expected 1 key after dedup, got %d", len(m))
	}
	if m["FOO"] != "second" {
		t.Errorf("expected last-write-wins, got %q", m["FOO"])
	}
}
