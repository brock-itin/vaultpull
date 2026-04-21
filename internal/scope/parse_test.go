package scope_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/scope"
)

func TestParseScopes_Basic(t *testing.T) {
	scopes, err := scope.ParseScopes([]string{"app:secret/myapp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scopes) != 1 {
		t.Fatalf("expected 1 scope, got %d", len(scopes))
	}
	if scopes[0].Name != "app" {
		t.Errorf("expected name=app, got %q", scopes[0].Name)
	}
	if scopes[0].Path != "secret/myapp" {
		t.Errorf("expected path=secret/myapp, got %q", scopes[0].Path)
	}
	if scopes[0].Prefix != "" {
		t.Errorf("expected empty prefix, got %q", scopes[0].Prefix)
	}
}

func TestParseScopes_WithPrefix(t *testing.T) {
	scopes, err := scope.ParseScopes([]string{"db:secret/database:DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scopes[0].Prefix != "DB_" {
		t.Errorf("expected prefix=DB_, got %q", scopes[0].Prefix)
	}
}

func TestParseScopes_Multiple(t *testing.T) {
	raw := []string{"app:secret/app", "infra:secret/infra:INFRA_"}
	scopes, err := scope.ParseScopes(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scopes) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(scopes))
	}
}

func TestParseScopes_MissingPath(t *testing.T) {
	_, err := scope.ParseScopes([]string{"onlynamewithoutcolon"})
	if err == nil {
		t.Fatal("expected error for missing path separator")
	}
}

func TestParseScopes_EmptyName(t *testing.T) {
	_, err := scope.ParseScopes([]string{":secret/app"})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestParseScopes_EmptyPath(t *testing.T) {
	_, err := scope.ParseScopes([]string{"app:"})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestParseScopes_Empty(t *testing.T) {
	scopes, err := scope.ParseScopes([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scopes) != 0 {
		t.Errorf("expected 0 scopes, got %d", len(scopes))
	}
}
