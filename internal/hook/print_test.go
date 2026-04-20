package hook_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/hook"
)

func TestPrint_Empty(t *testing.T) {
	var buf bytes.Buffer
	hook.Print(&buf, nil)
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty results, got %q", buf.String())
	}
}

func TestPrint_ShowsSuccess(t *testing.T) {
	var buf bytes.Buffer
	results := []hook.Result{
		{
			Hook:    hook.Hook{Stage: hook.StagePre, Command: "echo hi"},
			Output:  "hi",
			Err:     nil,
			Elapsed: 5 * time.Millisecond,
		},
	}
	hook.Print(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "ok") {
		t.Errorf("expected 'ok' in output, got: %s", out)
	}
	if !strings.Contains(out, "echo hi") {
		t.Errorf("expected command in output, got: %s", out)
	}
	if !strings.Contains(out, "hi") {
		t.Errorf("expected hook output in print, got: %s", out)
	}
}

func TestPrint_ShowsFailure(t *testing.T) {
	var buf bytes.Buffer
	results := []hook.Result{
		{
			Hook:    hook.Hook{Stage: hook.StagePost, Command: "false"},
			Err:     errors.New("exit status 1"),
			Elapsed: 2 * time.Millisecond,
		},
	}
	hook.Print(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "failed") {
		t.Errorf("expected 'failed' in output, got: %s", out)
	}
	if !strings.Contains(out, "exit status 1") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestPrint_MultipleResults(t *testing.T) {
	var buf bytes.Buffer
	results := []hook.Result{
		{Hook: hook.Hook{Stage: hook.StagePre, Command: "echo a"}, Output: "a"},
		{Hook: hook.Hook{Stage: hook.StagePre, Command: "echo b"}, Output: "b"},
	}
	hook.Print(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "echo a") || !strings.Contains(out, "echo b") {
		t.Errorf("expected both commands in output, got: %s", out)
	}
}
