package cli

import (
	"os"
	"testing"
)

func TestEnvOrDefault(t *testing.T) {
	os.Setenv("FILEFUNC_TEST_VAR", "custom")
	defer os.Unsetenv("FILEFUNC_TEST_VAR")

	if got := EnvOrDefault("FILEFUNC_TEST_VAR", "default"); got != "custom" {
		t.Errorf("got %q, want %q", got, "custom")
	}
	if got := EnvOrDefault("FILEFUNC_NONEXISTENT", "fallback"); got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func TestCheckProjectRoot(t *testing.T) {
	if err := CheckProjectRoot("/nonexistent/path"); err == nil {
		t.Error("expected error for nonexistent path")
	}

	tmp := t.TempDir()
	if err := CheckProjectRoot(tmp); err == nil {
		t.Error("expected error for missing go.mod")
	}

	os.WriteFile(tmp+"/go.mod", []byte("module test"), 0644)
	if err := CheckProjectRoot(tmp); err == nil {
		t.Error("expected error for missing codebook.yaml")
	}

	os.WriteFile(tmp+"/codebook.yaml", []byte(""), 0644)
	if err := CheckProjectRoot(tmp); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
