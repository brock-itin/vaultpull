package scope_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpull/internal/scope"
)

func TestResolve_NoPrefix(t *testing.T) {
	s := scope.Scope{Name: "app", Path: "secret/app"}
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}

	out, err := scope.Resolve(s, secrets, scope.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
}

func TestResolve_WithPrefix(t *testing.T) {
	s := scope.Scope{Name: "app", Path: "secret/app", Prefix: "APP_"}
	secrets := map[string]string{"KEY": "value"}

	out, err := scope.Resolve(s, secrets, scope.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_KEY"]; !ok {
		t.Error("expected key APP_KEY to exist")
	}
}

func TestResolve_StripPrefix(t *testing.T) {
	s := scope.Scope{Name: "app", Path: "secret/app", Prefix: "APP_"}
	secrets := map[string]string{"APP_KEY": "value"}
	opts := scope.Options{StripPrefix: true}

	out, err := scope.Resolve(s, secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["KEY"]; !ok {
		t.Error("expected stripped key KEY to exist")
	}
}

func TestResolve_EmptyPath(t *testing.T) {
	s := scope.Scope{Name: "bad", Path: ""}
	_, err := scope.Resolve(s, map[string]string{"K": "v"}, scope.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestResolveAll_MergesScopes(t *testing.T) {
	scopes := []scope.Scope{
		{Name: "base", Path: "secret/base", Prefix: "BASE_"},
		{Name: "override", Path: "secret/override", Prefix: "OVR_"},
	}
	data := map[string]map[string]string{
		"secret/base":     {"KEY": "base-val"},
		"secret/override": {"KEY": "ovr-val"},
	}
	fetch := func(path string) (map[string]string, error) {
		v, ok := data[path]
		if !ok {
			return nil, errors.New("not found")
		}
		return v, nil
	}

	out, err := scope.ResolveAll(scopes, fetch, scope.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["BASE_KEY"] != "base-val" {
		t.Errorf("expected BASE_KEY=base-val, got %q", out["BASE_KEY"])
	}
	if out["OVR_KEY"] != "ovr-val" {
		t.Errorf("expected OVR_KEY=ovr-val, got %q", out["OVR_KEY"])
	}
}

func TestResolveAll_FetchError(t *testing.T) {
	scopes := []scope.Scope{
		{Name: "bad", Path: "secret/missing"},
	}
	fetch := func(_ string) (map[string]string, error) {
		return nil, errors.New("vault unreachable")
	}
	_, err := scope.ResolveAll(scopes, fetch, scope.DefaultOptions())
	if err == nil {
		t.Fatal("expected error from failed fetch")
	}
}
