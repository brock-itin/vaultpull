package inject_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/inject"
)

func TestIntoProcess_SetsVariable(t *testing.T) {
	t.Setenv("VP_TEST_INJECT_KEY", "")
	os.Unsetenv("VP_TEST_INJECT_KEY")

	secrets := map[string]string{"VP_TEST_INJECT_KEY": "hello"}
	err := inject.IntoProcess(secrets, inject.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("VP_TEST_INJECT_KEY"); got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestIntoProcess_NoOverwrite(t *testing.T) {
	t.Setenv("VP_EXISTING", "original")

	secrets := map[string]string{"VP_EXISTING": "new"}
	opts := inject.DefaultOptions()
	opts.Overwrite = false

	if err := inject.IntoProcess(secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("VP_EXISTING"); got != "original" {
		t.Errorf("expected original value %q, got %q", "original", got)
	}
}

func TestIntoProcess_Overwrite(t *testing.T) {
	t.Setenv("VP_OVERWRITE", "old")

	secrets := map[string]string{"VP_OVERWRITE": "new"}
	opts := inject.DefaultOptions()
	opts.Overwrite = true

	if err := inject.IntoProcess(secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("VP_OVERWRITE"); got != "new" {
		t.Errorf("expected %q, got %q", "new", got)
	}
}

func TestIntoProcess_WithPrefix(t *testing.T) {
	os.Unsetenv("APP_SECRET_KEY")

	secrets := map[string]string{"SECRET_KEY": "abc123"}
	opts := inject.DefaultOptions()
	opts.Prefix = "APP_"

	if err := inject.IntoProcess(secrets, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("APP_SECRET_KEY"); got != "abc123" {
		t.Errorf("expected %q, got %q", "abc123", got)
	}
}

func TestIntoCommand_InjectsEnv(t *testing.T) {
	cmd := exec.Command("env")
	secrets := map[string]string{"VP_CMD_TEST": "injected"}

	if err := inject.IntoCommand(cmd, secrets, inject.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, e := range cmd.Env {
		if e == "VP_CMD_TEST=injected" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected VP_CMD_TEST=injected in cmd.Env")
	}
}

func TestIntoCommand_InheritsProcess(t *testing.T) {
	t.Setenv("VP_INHERIT_CHECK", "yes")

	cmd := exec.Command("env")
	if err := inject.IntoCommand(cmd, map[string]string{}, inject.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, e := range cmd.Env {
		if strings.HasPrefix(e, "VP_INHERIT_CHECK=") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected VP_INHERIT_CHECK in cmd.Env")
	}
}
