package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/vaultpull/internal/diff"
)

func TestPrint_Empty(t *testing.T) {
	var buf bytes.Buffer
	diff.Print(&buf, []diff.Change{})
	if !strings.Contains(buf.String(), "No changes") {
		t.Fatalf("expected no changes message, got: %s", buf.String())
	}
}

func TestPrint_ShowsAllTypes(t *testing.T) {
	changes := []diff.Change{
		{Key: "ADDED_KEY", Type: diff.Added, NewValue: "secret"},
		{Key: "UPDATED_KEY", Type: diff.Updated, OldValue: "old", NewValue: "new"},
		{Key: "DELETED_KEY", Type: diff.Deleted, OldValue: "gone"},
		{Key: "SAME_KEY", Type: diff.Unchanged, OldValue: "same", NewValue: "same"},
	}

	var buf bytes.Buffer
	diff.Print(&buf, changes)
	out := buf.String()

	for _, want := range []string{"ADDED_KEY", "UPDATED_KEY", "DELETED_KEY", "SAME_KEY", "Summary"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestPrint_MasksValues(t *testing.T) {
	changes := []diff.Change{
		{Key: "SECRET", Type: diff.Added, NewValue: "supersecret"},
	}

	var buf bytes.Buffer
	diff.Print(&buf, changes)
	if strings.Contains(buf.String(), "supersecret") {
		t.Fatal("print output should not contain raw secret values")
	}
}
