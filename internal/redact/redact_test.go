package redact_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/redact"
)

func TestMap_RedactsDefaultKeywords(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_TOKEN":   "tok123",
		"APP_NAME":    "myapp",
	}
	out := redact.Map(secrets, nil)
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD redacted, got %s", out["DB_PASSWORD"])
	}
	if out["API_TOKEN"] != "[REDACTED]" {
		t.Errorf("expected API_TOKEN redacted, got %s", out["API_TOKEN"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME preserved, got %s", out["APP_NAME"])
	}
}

func TestMap_CustomPlaceholder(t *testing.T) {
	secrets := map[string]string{"SECRET_KEY": "abc"}
	out := redact.Map(secrets, &redact.Options{Placeholder: "***"})
	if out["SECRET_KEY"] != "***" {
		t.Errorf("expected ***, got %s", out["SECRET_KEY"])
	}
}

func TestMap_CustomKeywords(t *testing.T) {
	secrets := map[string]string{
		"STRIPE_APIKEY": "sk_live_123",
		"DB_HOST":       "localhost",
	}
	opts := &redact.Options{Keywords: []string{"apikey"}}
	out := redact.Map(secrets, opts)
	if out["STRIPE_APIKEY"] != "[REDACTED]" {
		t.Errorf("expected STRIPE_APIKEY redacted")
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST preserved")
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"PASSWORD": "original"}
	redact.Map(secrets, nil)
	if secrets["PASSWORD"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"PRIVATE_KEY", true},
		{"APP_ENV", false},
		{"CREDENTIAL_STORE", true},
		{"PORT", false},
	}
	for _, tc := range cases {
		got := redact.IsSensitive(tc.key, nil)
		if got != tc.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}
