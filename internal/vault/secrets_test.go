package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestGetSecrets_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]string{"DB_PASS": "secret123"},
		},
	}
	ts := newTestServer(t, http.StatusOK, payload)
	defer ts.Close()

	c, _ := NewClient(ts.URL, "tok")
	secrets, err := c.GetSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["DB_PASS"] != "secret123" {
		t.Errorf("expected DB_PASS=secret123, got %q", secrets["DB_PASS"])
	}
}

func TestGetSecrets_NotFound(t *testing.T) {
	ts := newTestServer(t, http.StatusNotFound, nil)
	defer ts.Close()

	c, _ := NewClient(ts.URL, "tok")
	_, err := c.GetSecrets("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetSecrets_Forbidden(t *testing.T) {
	ts := newTestServer(t, http.StatusForbidden, nil)
	defer ts.Close()

	c, _ := NewClient(ts.URL, "badtoken")
	_, err := c.GetSecrets("secret/data/myapp")
	if err == nil {
		t.Fatal("expected error for 403, got nil")
	}
}
