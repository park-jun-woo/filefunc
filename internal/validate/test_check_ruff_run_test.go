//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckRuffRun — ruff 설치 시 린트 위반에 N4 violation 생성 확인
package validate

import (
	"os/exec"
	"testing"
)

func TestCheckRuffRun(t *testing.T) {
	if _, err := exec.LookPath("ruff"); err != nil {
		t.Skip("ruff not installed; skipping")
	}
	violations := CheckRuff([]string{"testdata/py_ruff_violation.py"})
	if len(violations) == 0 {
		t.Error("expected N4 violation for ruff lint issue, got none")
		return
	}
	if violations[0].Rule != "N4" {
		t.Errorf("expected rule N4, got %s", violations[0].Rule)
	}
}
