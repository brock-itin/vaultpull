package hook_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/hook"
)

func TestRun_PreHookExecuted(t *testing.T) {
	r := hook.New([]hook.Hook{
		{Stage: hook.StagePre, Command: "echo hello"},
	})
	results := r.Run(context.Background(), hook.StagePre)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].Output != "hello" {
		t.Errorf("expected output 'hello', got %q", results[0].Output)
	}
}

func TestRun_SkipsWrongStage(t *testing.T) {
	r := hook.New([]hook.Hook{
		{Stage: hook.StagePost, Command: "echo post"},
	})
	results := r.Run(context.Background(), hook.StagePre)
	if len(results) != 0 {
		t.Errorf("expected 0 results for wrong stage, got %d", len(results))
	}
}

func TestRun_CommandFailure(t *testing.T) {
	r := hook.New([]hook.Hook{
		{Stage: hook.StagePre, Command: "false"},
	})
	results := r.Run(context.Background(), hook.StagePre)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err == nil {
		t.Error("expected error for failing command")
	}
}

func TestRun_EmptyCommand(t *testing.T) {
	r := hook.New([]hook.Hook{
		{Stage: hook.StagePre, Command: ""},
	})
	results := r.Run(context.Background(), hook.StagePre)
	if results[0].Err == nil {
		t.Error("expected error for empty command")
	}
}

func TestRun_Timeout(t *testing.T) {
	r := hook.New([]hook.Hook{
		{Stage: hook.StagePre, Command: "sleep 10", Timeout: 50 * time.Millisecond},
	})
	results := r.Run(context.Background(), hook.StagePre)
	if results[0].Err == nil {
		t.Error("expected timeout error")
	}
}

func TestHasFailures_True(t *testing.T) {
	results := []hook.Result{
		{Err: nil},
		{Err: context.DeadlineExceeded},
	}
	if !hook.HasFailures(results) {
		t.Error("expected HasFailures to return true")
	}
}

func TestHasFailures_False(t *testing.T) {
	results := []hook.Result{
		{Err: nil},
		{Err: nil},
	}
	if hook.HasFailures(results) {
		t.Error("expected HasFailures to return false")
	}
}
