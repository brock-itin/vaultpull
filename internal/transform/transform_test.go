package transform_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/transform"
)

func TestApply_NoRule(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := transform.Apply(secrets, transform.Rule{})
	if len(result) != len(secrets) {
		t.Fatalf("expected %d keys, got %d", len(secrets), len(result))
	}
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", result["FOO"])
	}
}

func TestApply_AddPrefix(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	result := transform.Apply(secrets, transform.Rule{AddPrefix: "APP_"})
	if _, ok := result["APP_DB_HOST"]; !ok {
		t.Errorf("expected key APP_DB_HOST, got %v", result)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	secrets := map[string]string{"VAULT_SECRET_KEY": "value"}
	result := transform.Apply(secrets, transform.Rule{StripPrefix: "VAULT_"})
	if _, ok := result["SECRET_KEY"]; !ok {
		t.Errorf("expected key SECRET_KEY, got %v", result)
	}
}

func TestApply_Rename(t *testing.T) {
	secrets := map[string]string{"old_key": "value"}
	result := transform.Apply(secrets, transform.Rule{
		Rename: map[string]string{"old_key": "NEW_KEY"},
	})
	if result["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %v", result)
	}
	if _, ok := result["old_key"]; ok {
		t.Error("old_key should not exist after rename")
	}
}

func TestApply_PrefixThenRename(t *testing.T) {
	secrets := map[string]string{"host": "localhost"}
	result := transform.Apply(secrets, transform.Rule{
		AddPrefix: "DB_",
		Rename:    map[string]string{"DB_host": "DATABASE_HOST"},
	})
	if result["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %v", result)
	}
}

func TestUppercase(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost", "api_key": "secret"}
	result := transform.Uppercase(secrets)
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", result)
	}
	if result["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %v", result)
	}
}
