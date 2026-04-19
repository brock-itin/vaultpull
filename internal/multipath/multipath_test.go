package multipath_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpull/internal/multipath"
)

type mockResolver struct {
	data map[string]map[string]string
	errs map[string]error
}

func (m *mockResolver) GetSecrets(path string) (map[string]string, error) {
	if err, ok := m.errs[path]; ok {
		return nil, err
	}
	if d, ok := m.data[path]; ok {
		return d, nil
	}
	return map[string]string{}, nil
}

func TestMerge_Success(t *testing.T) {
	r := &mockResolver{
		data: map[string]map[string]string{
			"secret/a": {"FOO": "1", "BAR": "2"},
			"secret/b": {"BAZ": "3"},
		},
	}
	res := multipath.Merge(r, []string{"secret/a", "secret/b"})
	if res.HasErrors() {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if res.Secrets["FOO"] != "1" || res.Secrets["BAR"] != "2" || res.Secrets["BAZ"] != "3" {
		t.Errorf("unexpected secrets: %v", res.Secrets)
	}
}

func TestMerge_LaterPathWins(t *testing.T) {
	r := &mockResolver{
		data: map[string]map[string]string{
			"secret/a": {"KEY": "old"},
			"secret/b": {"KEY": "new"},
		},
	}
	res := multipath.Merge(r, []string{"secret/a", "secret/b"})
	if res.Secrets["KEY"] != "new" {
		t.Errorf("expected 'new', got %q", res.Secrets["KEY"])
	}
}

func TestMerge_PartialError(t *testing.T) {
	r := &mockResolver{
		data: map[string]map[string]string{"secret/a": {"FOO": "1"}},
		errs: map[string]error{"secret/b": errors.New("forbidden")},
	}
	res := multipath.Merge(r, []string{"secret/a", "secret/b"})
	if !res.HasErrors() {
		t.Fatal("expected errors")
	}
	if res.Secrets["FOO"] != "1" {
		t.Errorf("expected partial secrets to be present")
	}
}

func TestMerge_EmptyPath(t *testing.T) {
	r := &mockResolver{data: map[string]map[string]string{}}
	res := multipath.Merge(r, []string{""})
	if !res.HasErrors() {
		t.Fatal("expected error for empty path")
	}
}

func TestMergeStrict_Success(t *testing.T) {
	r := &mockResolver{
		data: map[string]map[string]string{
			"secret/a": {"A": "1"},
			"secret/b": {"B": "2"},
		},
	}
	got, err := multipath.MergeStrict(r, []string{"secret/a", "secret/b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["A"] != "1" || got["B"] != "2" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestMergeStrict_StopsOnError(t *testing.T) {
	r := &mockResolver{
		data: map[string]map[string]string{"secret/b": {"B": "2"}},
		errs: map[string]error{"secret/a": errors.New("not found")},
	}
	_, err := multipath.MergeStrict(r, []string{"secret/a", "secret/b"})
	if err == nil {
		t.Fatal("expected error")
	}
}
