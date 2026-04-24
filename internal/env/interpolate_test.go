package env

import (
	"testing"
)

func TestInterpolate_NoReferences(t *testing.T) {
	input := map[string]string{"KEY": "plainvalue"}
	out, err := Interpolate(input, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "plainvalue" {
		t.Errorf("expected plainvalue, got %q", out["KEY"])
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	input := map[string]string{
		"BASE": "/home/user",
		"PATH": "${BASE}/bin",
	}
	out, err := Interpolate(input, &InterpolateOptions{AllowMissing: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PATH"] != "/home/user/bin" {
		t.Errorf("expected /home/user/bin, got %q", out["PATH"])
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	input := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://$HOST/db",
	}
	out, err := Interpolate(input, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost/db" {
		t.Errorf("expected postgres://localhost/db, got %q", out["DSN"])
	}
}

func TestInterpolate_FromOptsEnv(t *testing.T) {
	input := map[string]string{"GREETING": "Hello ${NAME}"}
	opts := &InterpolateOptions{
		Env:          map[string]string{"NAME": "World"},
		AllowMissing: false,
	}
	out, err := Interpolate(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GREETING"] != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", out["GREETING"])
	}
}

func TestInterpolate_MissingAllowed(t *testing.T) {
	input := map[string]string{"VAL": "${UNDEFINED_XYZ_ABC}"}
	out, err := Interpolate(input, &InterpolateOptions{AllowMissing: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// unresolved reference is kept as-is
	if out["VAL"] != "${UNDEFINED_XYZ_ABC}" {
		t.Errorf("expected original placeholder, got %q", out["VAL"])
	}
}

func TestInterpolate_MissingDisallowed(t *testing.T) {
	input := map[string]string{"VAL": "${UNDEFINED_XYZ_ABC}"}
	_, err := Interpolate(input, &InterpolateOptions{AllowMissing: false})
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestInterpolate_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{
		"A": "hello",
		"B": "${A} world",
	}
	_, err := Interpolate(input, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input["B"] != "${A} world" {
		t.Errorf("input was mutated: got %q", input["B"])
	}
}
