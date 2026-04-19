package template

import (
	"strings"
	"testing"
)

func TestRender_BasicSubstitution(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "hunter2"}
	out, err := Render(`DB_PASS={{ index . "DB_PASS" }}`, secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "DB_PASS=hunter2" {
		t.Errorf("got %q", out)
	}
}

func TestRender_RequiredPresent(t *testing.T) {
	secrets := map[string]string{"API_KEY": "abc123"}
	tmpl := `KEY={{ required "API_KEY" . }}`
	out, err := Render(tmpl, secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "abc123") {
		t.Errorf("expected value in output, got %q", out)
	}
}

func TestRender_RequiredMissing(t *testing.T) {
	secrets := map[string]string{}
	tmpl := `KEY={{ required "MISSING" . }}`
	_, err := Render(tmpl, secrets, nil)
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestRender_DefaultFallback(t *testing.T) {
	secrets := map[string]string{}
	tmpl := `LEVEL={{ default "info" "LOG_LEVEL" . }}`
	out, err := Render(tmpl, secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "info") {
		t.Errorf("expected default value, got %q", out)
	}
}

func TestRender_DefaultOverridden(t *testing.T) {
	secrets := map[string]string{"LOG_LEVEL": "debug"}
	tmpl := `LEVEL={{ default "info" "LOG_LEVEL" . }}`
	out, err := Render(tmpl, secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "debug") {
		t.Errorf("expected overridden value, got %q", out)
	}
}

func TestRender_CustomDelimiters(t *testing.T) {
	secrets := map[string]string{"X": "42"}
	tmpl := `X=<< index . "X" >>`
	out, err := Render(tmpl, secrets, &Options{LeftDelim: "<<", RightDelim: ">>" })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "X=42" {
		t.Errorf("got %q", out)
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	_, err := Render(`{{ .Unclosed`, nil, nil)
	if err == nil {
		t.Fatal("expected parse error")
	}
}
