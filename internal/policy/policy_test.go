package policy_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/policy"
)

// helper to build a simple secrets map
func makeSecrets(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestCheck_NoRules(t *testing.T) {
	secrets := makeSecrets("DB_PASSWORD", "secret", "API_KEY", "abc123")
	result := policy.Check(secrets, policy.Options{})
	if policy.HasViolations(result) {
		t.Errorf("expected no violations with no rules, got %d", len(result))
	}
}

func TestCheck_RequiredKeyPresent(t *testing.T) {
	secrets := makeSecrets("DB_PASSWORD", "hunter2")
	opts := policy.Options{
		Required: []string{"DB_PASSWORD"},
	}
	result := policy.Check(secrets, opts)
	if policy.HasViolations(result) {
		t.Errorf("expected no violations, got: %v", result)
	}
}

func TestCheck_RequiredKeyMissing(t *testing.T) {
	secrets := makeSecrets("API_KEY", "abc")
	opts := policy.Options{
		Required: []string{"DB_PASSWORD", "DB_USER"},
	}
	result := policy.Check(secrets, opts)
	if !policy.HasViolations(result) {
		t.Fatal("expected violations for missing required keys")
	}
	if len(result) != 2 {
		t.Errorf("expected 2 violations, got %d", len(result))
	}
}

func TestCheck_ForbiddenKeyAbsent(t *testing.T) {
	secrets := makeSecrets("DB_PASSWORD", "safe")
	opts := policy.Options{
		Forbidden: []string{"DEBUG_TOKEN"},
	}
	result := policy.Check(secrets, opts)
	if policy.HasViolations(result) {
		t.Errorf("expected no violations, got: %v", result)
	}
}

func TestCheck_ForbiddenKeyPresent(t *testing.T) {
	secrets := makeSecrets("DB_PASSWORD", "safe", "DEBUG_TOKEN", "oops")
	opts := policy.Options{
		Forbidden: []string{"DEBUG_TOKEN"},
	}
	result := policy.Check(secrets, opts)
	if !policy.HasViolations(result) {
		t.Fatal("expected violation for forbidden key")
	}
	if len(result) != 1 {
		t.Errorf("expected 1 violation, got %d", len(result))
	}
}

func TestCheck_MinLengthPassed(t *testing.T) {
	secrets := makeSecrets("API_SECRET", "longenoughvalue")
	opts := policy.Options{
		MinLength: map[string]int{"API_SECRET": 8},
	}
	result := policy.Check(secrets, opts)
	if policy.HasViolations(result) {
		t.Errorf("expected no violations, got: %v", result)
	}
}

func TestCheck_MinLengthFailed(t *testing.T) {
	secrets := makeSecrets("API_SECRET", "short")
	opts := policy.Options{
		MinLength: map[string]int{"API_SECRET": 16},
	}
	result := policy.Check(secrets, opts)
	if !policy.HasViolations(result) {
		t.Fatal("expected violation for short value")
	}
}

func TestCheck_MinLengthKeyMissing(t *testing.T) {
	// key not present in secrets — should not panic, no violation for length
	secrets := makeSecrets("OTHER_KEY", "value")
	opts := policy.Options{
		MinLength: map[string]int{"API_SECRET": 8},
	}
	// missing key is only a violation if also listed in Required
	result := policy.Check(secrets, opts)
	if policy.HasViolations(result) {
		t.Errorf("expected no violations, got: %v", result)
	}
}

func TestCheck_CombinedRules(t *testing.T) {
	secrets := makeSecrets(
		"DB_PASSWORD", "ok",
		"DEBUG_TOKEN", "present",
	)
	opts := policy.Options{
		Required:  []string{"DB_PASSWORD", "API_KEY"},
		Forbidden: []string{"DEBUG_TOKEN"},
		MinLength: map[string]int{"DB_PASSWORD": 10},
	}
	result := policy.Check(secrets, opts)
	if !policy.HasViolations(result) {
		t.Fatal("expected violations")
	}
	// API_KEY missing (1) + DEBUG_TOKEN present (1) + DB_PASSWORD too short (1) = 3
	if len(result) != 3 {
		t.Errorf("expected 3 violations, got %d: %v", len(result), result)
	}
}

func TestHasViolations_Empty(t *testing.T) {
	if policy.HasViolations(nil) {
		t.Error("nil slice should not have violations")
	}
}
