package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func setupEnv(t *testing.T, addr, token, path string) {
	t.Helper()
	t.Setenv("VAULT_ADDR", addr)
	t.Setenv("VAULT_TOKEN", token)
	t.Setenv("VAULT_PATH", path)
}

func newRunTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]string{"APP_KEY": "secret123"},
			},
		}
		json.NewEncoder(w).Encode(payload)
	}))
}

func TestRun_Success(t *testing.T) {
	srv := newRunTestServer(t)
	defer srv.Close()

	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	setupEnv(t, srv.URL, "test-token", "secret/data/app")

	if err := Run(output, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if len(data) == 0 {
		t.Error("output file is empty")
	}
}

func TestRun_MissingToken(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://localhost:8200")
	t.Setenv("VAULT_TOKEN", "")
	t.Setenv("VAULT_PATH", "secret/data/app")

	if err := Run("", false); err == nil {
		t.Error("expected error for missing token")
	}
}
