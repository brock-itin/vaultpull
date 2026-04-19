package mask_test

import (
	"testing"

	"github.com/example/vaultpull/internal/mask"
)

func TestValue_NilOpts(t *testing.T) {
	if got := mask.Value("supersecret", nil); got != "****" {
		t.Fatalf("expected ****, got %q", got)
	}
}

func TestValue_EmptyString(t *testing.T) {
	if got := mask.Value("", nil); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestValue_ShowPrefix(t *testing.T) {
	opts := &mask.Options{ShowPrefix: 3}
	got := mask.Value("abcdef", opts)
	if got != "abc****" {
		t.Fatalf("expected abc****, got %q", got)
	}
}

func TestValue_PrefixLongerThanSecret(t *testing.T) {
	opts := &mask.Options{ShowPrefix: 20}
	got := mask.Value("hi", opts)
	if got != "****" {
		t.Fatalf("expected ****, got %q", got)
	}
}

func TestValue_ZeroPrefix(t *testing.T) {
	opts := &mask.Options{ShowPrefix: 0}
	if got := mask.Value("secret", opts); got != "****" {
		t.Fatalf("expected ****, got %q", got)
	}
}

func TestMap_MasksAllValues(t *testing.T) {
	input := map[string]string{
		"KEY_A": "alpha",
		"KEY_B": "beta",
	}
	out := mask.Map(input, nil)
	for k, v := range out {
		if v != "****" {
			t.Errorf("key %s: expected ****, got %q", k, v)
		}
	}
}

func TestMap_PreservesKeys(t *testing.T) {
	input := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := mask.Map(input, nil)
	for k := range input {
		if _, ok := out[k]; !ok {
			t.Errorf("key %s missing from masked map", k)
		}
	}
}
