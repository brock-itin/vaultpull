package prompt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/prompt"
)

func TestConfirm_Yes(t *testing.T) {
	for _, input := range []string{"y\n", "Y\n", "yes\n", "YES\n", "Yes\n"} {
		p := prompt.New(strings.NewReader(input), &bytes.Buffer{})
		ok, err := p.Confirm("Continue?")
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if !ok {
			t.Errorf("input %q: expected true, got false", input)
		}
	}
}

func TestConfirm_No(t *testing.T) {
	for _, input := range []string{"n\n", "no\n", "\n", "maybe\n"} {
		p := prompt.New(strings.NewReader(input), &bytes.Buffer{})
		ok, err := p.Confirm("Continue?")
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if ok {
			t.Errorf("input %q: expected false, got true", input)
		}
	}
}

func TestConfirm_EOF(t *testing.T) {
	p := prompt.New(strings.NewReader(""), &bytes.Buffer{})
	ok, err := p.Confirm("Continue?")
	if err != nil {
		t.Fatalf("unexpected error on EOF: %v", err)
	}
	if ok {
		t.Error("expected false on EOF")
	}
}

func TestConfirm_WritesQuestion(t *testing.T) {
	var out bytes.Buffer
	p := prompt.New(strings.NewReader("y\n"), &out)
	_, _ = p.Confirm("Apply changes?")
	if !strings.Contains(out.String(), "Apply changes?") {
		t.Errorf("expected question in output, got: %q", out.String())
	}
}

func TestAutoConfirm_AlwaysTrue(t *testing.T) {
	a := prompt.AutoConfirm{}
	ok, err := a.Confirm("anything")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected AutoConfirm to return true")
	}
}
