package cmd

import "testing"

func TestParseFlags_Defaults(t *testing.T) {
	opts, err := ParseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Output != ".env" {
		t.Errorf("expected .env, got %s", opts.Output)
	}
	if opts.Overwrite {
		t.Error("expected overwrite=false by default")
	}
}

func TestParseFlags_CustomOutput(t *testing.T) {
	opts, err := ParseFlags([]string{"-output", "secrets/.env", "-overwrite"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Output != "secrets/.env" {
		t.Errorf("expected secrets/.env, got %s", opts.Output)
	}
	if !opts.Overwrite {
		t.Error("expected overwrite=true")
	}
}

func TestParseFlags_InvalidFlag(t *testing.T) {
	_, err := ParseFlags([]string{"-unknown"})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}
