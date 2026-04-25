package healthcheck_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpull/internal/healthcheck"
)

func newHealthServer(sysStatus int, tokenStatus int) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/sys/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Vault-Version", "1.15.0")
		w.WriteHeader(sysStatus)
	})
	mux.HandleFunc("/v1/auth/token/lookup-self", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(tokenStatus)
	})
	return httptest.NewServer(mux)
}

func TestCheck_ReachableAndTokenValid(t *testing.T) {
	srv := newHealthServer(http.StatusOK, http.StatusOK)
	defer srv.Close()

	s := healthcheck.Check(context.Background(), srv.URL, "root", healthcheck.DefaultOptions())

	if !s.Reachable {
		t.Error("expected Reachable=true")
	}
	if !s.TokenValid {
		t.Error("expected TokenValid=true")
	}
	if s.Error != nil {
		t.Errorf("unexpected error: %v", s.Error)
	}
	if s.VaultVersion != "1.15.0" {
		t.Errorf("expected version 1.15.0, got %q", s.VaultVersion)
	}
}

func TestCheck_InvalidToken(t *testing.T) {
	srv := newHealthServer(http.StatusOK, http.StatusForbidden)
	defer srv.Close()

	s := healthcheck.Check(context.Background(), srv.URL, "bad-token", healthcheck.DefaultOptions())

	if !s.Reachable {
		t.Error("expected Reachable=true")
	}
	if s.TokenValid {
		t.Error("expected TokenValid=false")
	}
	if s.Error == nil {
		t.Error("expected an error for invalid token")
	}
}

func TestCheck_Unreachable(t *testing.T) {
	s := healthcheck.Check(context.Background(), "http://127.0.0.1:19999", "root", healthcheck.DefaultOptions())

	if s.Reachable {
		t.Error("expected Reachable=false for unreachable server")
	}
	if s.Error == nil {
		t.Error("expected an error for unreachable server")
	}
}

func TestOK_BothTrue(t *testing.T) {
	s := healthcheck.Status{Reachable: true, TokenValid: true}
	if !healthcheck.OK(s) {
		t.Error("expected OK=true")
	}
}

func TestOK_OnlyReachable(t *testing.T) {
	s := healthcheck.Status{Reachable: true, TokenValid: false}
	if healthcheck.OK(s) {
		t.Error("expected OK=false when token invalid")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := healthcheck.DefaultOptions()
	if opts.Timeout <= 0 {
		t.Error("expected positive default timeout")
	}
}
