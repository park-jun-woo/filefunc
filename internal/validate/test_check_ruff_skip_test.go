//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckRuffNotInstalled — ruff 미설치 시 N4 ERROR 반환 확인
package validate

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCheckRuffNotInstalled(t *testing.T) {
	if _, err := exec.LookPath("ruff"); err == nil {
		t.Skip("ruff is installed; skip not-installed test")
	}
	violations := CheckRuff([]string{"testdata/py_ruff_violation.py"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "ruff not installed") {
		t.Errorf("expected 'ruff not installed' message, got %q", violations[0].Message)
	}
}
