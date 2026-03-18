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

func TestFindGoMod(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(tmp+"/go.mod", []byte("module test"), 0644)
	os.MkdirAll(tmp+"/sub/deep", 0755)

	got := FindGoMod(tmp + "/sub/deep")
	if got != tmp+"/go.mod" {
		t.Errorf("FindGoMod = %q, want %q", got, tmp+"/go.mod")
	}
}

func TestFindGoMod_NotFound(t *testing.T) {
	got := FindGoMod("/nonexistent/deep/path")
	if got != "go.mod" {
		t.Errorf("FindGoMod = %q, want %q", got, "go.mod")
	}
}
