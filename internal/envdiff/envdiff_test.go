package envdiff_test

import (
	"bytes"
	"testing"

	"github.com/your-org/vaultpull/internal/envdiff"
)

func TestCompare_Added(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"FOO": "bar"}
	r := envdiff.Compare(old, new)
	if len(r.Entries) != 1 || r.Entries[0].Change != envdiff.Added {
		t.Fatalf("expected 1 Added entry, got %+v", r.Entries)
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{}
	r := envdiff.Compare(old, new)
	if len(r.Entries) != 1 || r.Entries[0].Change != envdiff.Removed {
		t.Fatalf("expected 1 Removed entry, got %+v", r.Entries)
	}
}

func TestCompare_Changed(t *testing.T) {
	old := map[string]string{"FOO": "old"}
	new := map[string]string{"FOO": "new"}
	r := envdiff.Compare(old, new)
	if len(r.Entries) != 1 || r.Entries[0].Change != envdiff.Changed {
		t.Fatalf("expected 1 Changed entry, got %+v", r.Entries)
	}
	if r.Entries[0].Old != "old" || r.Entries[0].New != "new" {
		t.Fatalf("unexpected values: %+v", r.Entries[0])
	}
}

func TestCompare_Same(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{"FOO": "bar"}
	r := envdiff.Compare(old, new)
	if r.HasChanges() {
		t.Fatal("expected no changes")
	}
	if r.Entries[0].Change != envdiff.Same {
		t.Fatalf("expected Same, got %s", r.Entries[0].Change)
	}
}

func TestCompare_Empty(t *testing.T) {
	r := envdiff.Compare(map[string]string{}, map[string]string{})
	if r.HasChanges() {
		t.Fatal("expected no changes for empty maps")
	}
	if len(r.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(r.Entries))
	}
}

func TestSummary(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "changed", "C": "3"}
	r := envdiff.Compare(old, new)
	s := r.Summary()
	if s[envdiff.Changed] != 1 {
		t.Errorf("expected 1 changed, got %d", s[envdiff.Changed])
	}
	if s[envdiff.Removed] != 1 {
		t.Errorf("expected 1 removed, got %d", s[envdiff.Removed])
	}
	if s[envdiff.Added] != 1 {
		t.Errorf("expected 1 added, got %d", s[envdiff.Added])
	}
}

func TestPrint_MasksValues(t *testing.T) {
	old := map[string]string{"SECRET": "old-val"}
	new := map[string]string{"SECRET": "new-val"}
	r := envdiff.Compare(old, new)
	var buf bytes.Buffer
	envdiff.Print(&buf, r, true)
	if bytes.Contains(buf.Bytes(), []byte("old-val")) || bytes.Contains(buf.Bytes(), []byte("new-val")) {
		t.Fatalf("expected values to be masked, got: %s", buf.String())
	}
}

func TestPrint_NoChanges(t *testing.T) {
	r := envdiff.Compare(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	var buf bytes.Buffer
	envdiff.Print(&buf, r, false)
	if !bytes.Contains(buf.Bytes(), []byte("No changes")) {
		t.Fatalf("expected no-changes message, got: %s", buf.String())
	}
}
