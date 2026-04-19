package validate_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/validate"
)

func TestCheck_RequiredPresent(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "s3cr3t", "API_KEY": "abc"}
	issues := validate.Check(secrets, validate.Options{Required: []string{"DB_PASS", "API_KEY"}})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestCheck_RequiredMissing(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "s3cr3t"}
	issues := validate.Check(secrets, validate.Options{Required: []string{"DB_PASS", "API_KEY"}})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY issue, got %s", issues[0].Key)
	}
}

func TestCheck_RequiredEmpty(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "   "}
	issues := validate.Check(secrets, validate.Options{Required: []string{"DB_PASS"}})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue for blank value, got %d", len(issues))
	}
}

func TestCheck_PlaceholderDetected(t *testing.T) {
	secrets := map[string]string{"TOKEN": "CHANGEME", "HOST": "localhost"}
	issues := validate.Check(secrets, validate.Options{ForbidPlaceholders: true})
	if len(issues) != 1 {
		t.Fatalf("expected 1 placeholder issue, got %d", len(issues))
	}
	if issues[0].Key != "TOKEN" {
		t.Errorf("expected TOKEN, got %s", issues[0].Key)
	}
}

func TestCheck_PlaceholderCaseInsensitive(t *testing.T) {
	secrets := map[string]string{"SECRET": "changeme"}
	issues := validate.Check(secrets, validate.Options{ForbidPlaceholders: true})
	if len(issues) == 0 {
		t.Fatal("expected placeholder issue for lowercase changeme")
	}
}

func TestCheck_NoOptions(t *testing.T) {
	secrets := map[string]string{"A": "", "B": "CHANGEME"}
	issues := validate.Check(secrets, validate.Options{})
	if len(issues) != 0 {
		t.Fatalf("expected no issues with empty options, got %v", issues)
	}
}

func TestHasIssues(t *testing.T) {
	secrets := map[string]string{}
	if !validate.HasIssues(secrets, validate.Options{Required: []string{"MUST_EXIST"}}) {
		t.Error("expected HasIssues to return true")
	}
}
