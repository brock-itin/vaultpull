package diff_test

import (
	"testing"

	"github.com/yourorg/vaultpull/internal/diff"
)

func TestCompare_Added(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"FOO": "bar"}

	changes := diff.Compare(existing, incoming)
	if len(changes) != 1 || changes[0].Type != diff.Added {
		t.Fatalf("expected 1 Added change, got %+v", changes)
	}
}

func TestCompare_Updated(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}

	changes := diff.Compare(existing, incoming)
	if len(changes) != 1 || changes[0].Type != diff.Updated {
		t.Fatalf("expected 1 Updated change, got %+v", changes)
	}
	if changes[0].OldValue != "old" || changes[0].NewValue != "new" {
		t.Fatalf("unexpected values: %+v", changes[0])
	}
}

func TestCompare_Deleted(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{}

	changes := diff.Compare(existing, incoming)
	if len(changes) != 1 || changes[0].Type != diff.Deleted {
		t.Fatalf("expected 1 Deleted change, got %+v", changes)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar"}

	changes := diff.Compare(existing, incoming)
	if len(changes) != 1 || changes[0].Type != diff.Unchanged {
		t.Fatalf("expected 1 Unchanged change, got %+v", changes)
	}
}

func TestSummary(t *testing.T) {
	changes := []diff.Change{
		{Type: diff.Added},
		{Type: diff.Added},
		{Type: diff.Updated},
		{Type: diff.Deleted},
	}
	s := diff.Summary(changes)
	if s[diff.Added] != 2 || s[diff.Updated] != 1 || s[diff.Deleted] != 1 {
		t.Fatalf("unexpected summary: %+v", s)
	}
}
