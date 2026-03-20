//ff:func feature=cli type=util control=sequence
//ff:what test: TestEnvOrDefault
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
