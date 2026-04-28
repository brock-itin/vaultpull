package chain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/chain"
)

func identity(m map[string]string) (map[string]string, error) {
	return m, nil
}

func uppercaseValues(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = strings.ToUpper(v)
	}
	return out, nil
}

func failStep(_ map[string]string) (map[string]string, error) {
	return nil, errors.New("step failed")
}

func nilStep(_ map[string]string) (map[string]string, error) {
	return nil, nil
}

func TestRun_EmptyPipeline(t *testing.T) {
	input := map[string]string{"KEY": "val"}
	res, err := chain.New().Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["KEY"] != "val" {
		t.Errorf("expected val, got %s", res.Secrets["KEY"])
	}
}

func TestRun_SingleStep(t *testing.T) {
	input := map[string]string{"KEY": "hello"}
	res, err := chain.New().Add("upper", uppercaseValues).Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["KEY"] != "HELLO" {
		t.Errorf("expected HELLO, got %s", res.Secrets["KEY"])
	}
	if len(res.Applied) != 1 || res.Applied[0] != "upper" {
		t.Errorf("expected applied=[upper], got %v", res.Applied)
	}
}

func TestRun_MultipleSteps(t *testing.T) {
	input := map[string]string{"A": "x"}
	res, err := chain.New().
		Add("id", identity).
		Add("upper", uppercaseValues).
		Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["A"] != "X" {
		t.Errorf("expected X, got %s", res.Secrets["A"])
	}
	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied steps, got %d", len(res.Applied))
	}
}

func TestRun_StepError_Halts(t *testing.T) {
	input := map[string]string{"K": "v"}
	_, err := chain.New().
		Add("ok", identity).
		Add("bad", failStep).
		Add("never", identity).
		Run(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "bad") {
		t.Errorf("expected step name in error, got: %v", err)
	}
}

func TestRun_NilReturnSkipsStep(t *testing.T) {
	input := map[string]string{"K": "v"}
	res, err := chain.New().
		Add("skip", nilStep).
		Add("upper", uppercaseValues).
		Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "skip" {
		t.Errorf("expected skipped=[skip], got %v", res.Skipped)
	}
	if res.Secrets["K"] != "V" {
		t.Errorf("expected V, got %s", res.Secrets["K"])
	}
}

func TestRun_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"K": "original"}
	_, _ = chain.New().Add("upper", uppercaseValues).Run(input)
	if input["K"] != "original" {
		t.Errorf("input was mutated, got %s", input["K"])
	}
}
