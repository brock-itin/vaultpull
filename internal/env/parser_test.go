package env

import (
	"strings"
	"testing"
)

func TestParse_BasicKeyValue(t *testing.T) {
	input := "FOO=bar\nBAZ=qux\n"
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "FOO" || entries[0].Value != "bar" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestParse_QuotedValues(t *testing.T) {
	input := `DB_URL="postgres://localhost/mydb"` + "\n" +
		`SECRET='top secret'` + "\n"
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "postgres://localhost/mydb" {
		t.Errorf("expected unquoted value, got %q", entries[0].Value)
	}
	if entries[1].Value != "top secret" {
		t.Errorf("expected unquoted value, got %q", entries[1].Value)
	}
}

func TestParse_CommentsAndBlanks(t *testing.T) {
	input := "# comment\n\nFOO=bar\n"
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "" || entries[1].Key != "" {
		t.Error("comment and blank lines should have empty Key")
	}
	if entries[2].Key != "FOO" {
		t.Errorf("expected FOO, got %q", entries[2].Key)
	}
}

func TestParse_MissingEquals(t *testing.T) {
	input := "INVALID_LINE\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParse_EmptyKey(t *testing.T) {
	input := "=value\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestToMap(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "", Raw: "# comment"},
		{Key: "B", Value: "2"},
	}
	m := ToMap(entries)
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected map: %v", m)
	}
	if _, ok := m[""]; ok {
		t.Error("empty key should not appear in map")
	}
}
