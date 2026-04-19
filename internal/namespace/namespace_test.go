package namespace_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/namespace"
)

func TestResolve_Basic(t *testing.T) {
	r := namespace.New(namespace.Options{
		Base: "secret/data",
		Env:  "prod",
	})
	got, err := r.Resolve("database")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "secret/data/prod/database"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolve_WithTeam(t *testing.T) {
	r := namespace.New(namespace.Options{
		Base: "secret/data",
		Env:  "staging",
		Team: "payments",
	})
	got, err := r.Resolve("stripe")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "secret/data/staging/payments/stripe"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolve_EmptyName(t *testing.T) {
	r := namespace.New(namespace.Options{Base: "secret/data", Env: "prod"})
	_, err := r.Resolve("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestResolve_EmptyBase(t *testing.T) {
	r := namespace.New(namespace.Options{Env: "prod"})
	_, err := r.Resolve("db")
	if err == nil {
		t.Fatal("expected error for empty base")
	}
}

func TestResolve_EmptyEnv(t *testing.T) {
	r := namespace.New(namespace.Options{Base: "secret/data"})
	_, err := r.Resolve("db")
	if err == nil {
		t.Fatal("expected error for empty env")
	}
}

func TestResolveAll_Success(t *testing.T) {
	r := namespace.New(namespace.Options{Base: "secret/data", Env: "dev"})
	paths, err := r.ResolveAll([]string{"alpha", "beta"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
	if paths[0] != "secret/data/dev/alpha" {
		t.Errorf("unexpected path[0]: %q", paths[0])
	}
	if paths[1] != "secret/data/dev/beta" {
		t.Errorf("unexpected path[1]: %q", paths[1])
	}
}

func TestResolveAll_ErrorPropagates(t *testing.T) {
	r := namespace.New(namespace.Options{Base: "secret/data", Env: "dev"})
	_, err := r.ResolveAll([]string{"ok", ""})
	if err == nil {
		t.Fatal("expected error for empty name in list")
	}
}
