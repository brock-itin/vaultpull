package truncate_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/truncate"
)

func TestValue_ShortString(t *testing.T) {
	out := truncate.Value("hello", nil)
	if out != "hello" {
		t.Fatalf("expected 'hello', got %q", out)
	}
}

func TestValue_ExactLength(t *testing.T) {
	s := strings.Repeat("a", 64)
	out := truncate.Value(s, nil)
	if out != s {
		t.Fatal("expected string unchanged at exact max length")
	}
}

func TestValue_Truncated(t *testing.T) {
	s := strings.Repeat("x", 100)
	out := truncate.Value(s, nil)
	if len(out) != 64+3 {
		t.Fatalf("expected length %d, got %d", 64+3, len(out))
	}
	if !strings.HasSuffix(out, "...") {
		t.Fatal("expected suffix '...'")
	}
}

func TestValue_CustomOpts(t *testing.T) {
	opts := &truncate.Options{MaxLen: 5, Suffix: "~"}
	out := truncate.Value("abcdefgh", opts)
	if out != "abcde~" {
		t.Fatalf("expected 'abcde~', got %q", out)
	}
}

func TestValue_EmptyString(t *testing.T) {
	out := truncate.Value("", nil)
	if out != "" {
		t.Fatalf("expected empty string, got %q", out)
	}
}

func TestMap_TruncatesValues(t *testing.T) {
	m := map[string]string{
		"SHORT": "hi",
		"LONG":  strings.Repeat("z", 100),
	}
	out := truncate.Map(m, nil)
	if out["SHORT"] != "hi" {
		t.Errorf("SHORT should be unchanged")
	}
	if !strings.HasSuffix(out["LONG"], "...") {
		t.Errorf("LONG should be truncated with suffix")
	}
	if len(out["LONG"]) != 67 {
		t.Errorf("LONG expected length 67, got %d", len(out["LONG"]))
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	original := strings.Repeat("a", 100)
	m := map[string]string{"KEY": original}
	truncate.Map(m, nil)
	if m["KEY"] != original {
		t.Fatal("input map was mutated")
	}
}
