package label_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/label"
)

func TestParse_Valid(t *testing.T) {
	lbls, err := label.Parse([]string{"env=prod", "team=backend"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lbls["env"] != "prod" || lbls["team"] != "backend" {
		t.Errorf("unexpected labels: %v", lbls)
	}
}

func TestParse_MissingEquals(t *testing.T) {
	_, err := label.Parse([]string{"noequalssign"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParse_EmptyKey(t *testing.T) {
	_, err := label.Parse([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestMatch_AllRequired(t *testing.T) {
	lbls := label.Labels{"env": "prod", "team": "backend"}
	opts := label.Options{Required: label.Labels{"env": "prod"}}
	if !label.Match(lbls, opts) {
		t.Error("expected match")
	}
}

func TestMatch_RequiredMissing(t *testing.T) {
	lbls := label.Labels{"env": "staging"}
	opts := label.Options{Required: label.Labels{"env": "prod"}}
	if label.Match(lbls, opts) {
		t.Error("expected no match")
	}
}

func TestMatch_ForbiddenPresent(t *testing.T) {
	lbls := label.Labels{"env": "prod", "deprecated": "true"}
	opts := label.Options{Forbidden: []string{"deprecated"}}
	if label.Match(lbls, opts) {
		t.Error("expected no match due to forbidden key")
	}
}

func TestMatch_ForbiddenAbsent(t *testing.T) {
	lbls := label.Labels{"env": "prod"}
	opts := label.Options{Forbidden: []string{"deprecated"}}
	if !label.Match(lbls, opts) {
		t.Error("expected match when forbidden key absent")
	}
}

func TestFilter_SelectsMatchingKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS":  "secret1",
		"API_KEY":  "secret2",
		"OLD_PASS": "secret3",
	}
	secretLabels := map[string]label.Labels{
		"DB_PASS":  {"env": "prod"},
		"API_KEY":  {"env": "prod"},
		"OLD_PASS": {"env": "staging"},
	}
	opts := label.Options{Required: label.Labels{"env": "prod"}}
	out := label.Filter(secrets, secretLabels, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	if _, ok := out["OLD_PASS"]; ok {
		t.Error("OLD_PASS should have been filtered out")
	}
}

func TestFilter_EmptyOpts_ReturnsAll(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	secretLabels := map[string]label.Labels{}
	out := label.Filter(secrets, secretLabels, label.DefaultOptions())
	if len(out) != 2 {
		t.Errorf("expected all secrets, got %d", len(out))
	}
}
