package sanitize_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/sanitize"
)

func TestKey_Uppercase(t *testing.T) {
	opts := sanitize.Options{UppercaseKeys: true}
	got := sanitize.Key("my_secret", opts)
	if got != "MY_SECRET" {
		t.Errorf("expected MY_SECRET, got %s", got)
	}
}

func TestKey_ReplaceHyphens(t *testing.T) {
	opts := sanitize.Options{ReplaceHyphens: true}
	got := sanitize.Key("my-secret-key", opts)
	if got != "my_secret_key" {
		t.Errorf("expected my_secret_key, got %s", got)
	}
}

func TestKey_StripInvalidChars(t *testing.T) {
	opts := sanitize.Options{StripInvalidChars: true}
	got := sanitize.Key("my.secret@key!", opts)
	if got != "mysecretkey" {
		t.Errorf("expected mysecretkey, got %s", got)
	}
}

func TestKey_StripLeadingDigits(t *testing.T) {
	opts := sanitize.DefaultOptions()
	got := sanitize.Key("123secret", opts)
	if got != "SECRET" {
		t.Errorf("expected SECRET, got %s", got)
	}
}

func TestKey_AllOptions(t *testing.T) {
	opts := sanitize.DefaultOptions()
	got := sanitize.Key("my-secret.key", opts)
	if got != "MY_SECRETKEY" {
		t.Errorf("expected MY_SECRETKEY, got %s", got)
	}
}

func TestKey_EmptyResult_AfterStrip(t *testing.T) {
	opts := sanitize.Options{StripInvalidChars: true}
	got := sanitize.Key("!!!@@@", opts)
	if got != "" {
		t.Errorf("expected empty string, got %s", got)
	}
}

func TestMap_SanitizesKeys(t *testing.T) {
	input := map[string]string{
		"my-key":   "value1",
		"other_key": "value2",
	}
	opts := sanitize.DefaultOptions()
	out := sanitize.Map(input, opts)

	if _, ok := out["MY_KEY"]; !ok {
		t.Error("expected MY_KEY in output")
	}
	if _, ok := out["OTHER_KEY"]; !ok {
		t.Error("expected OTHER_KEY in output")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestMap_DropsEmptyKeys(t *testing.T) {
	input := map[string]string{
		"!!!": "value",
		"valid": "ok",
	}
	opts := sanitize.DefaultOptions()
	out := sanitize.Map(input, opts)

	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
	if out["VALID"] != "ok" {
		t.Errorf("expected VALID=ok, got %s", out["VALID"])
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := sanitize.DefaultOptions()
	if !opts.UppercaseKeys || !opts.ReplaceHyphens || !opts.StripInvalidChars {
		t.Error("expected all default options to be true")
	}
}
