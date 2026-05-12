//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckRuffSkip — ruff 미설치 시 빈 슬라이스 반환 확인
package validate

import (
	"os/exec"
	"testing"
)

func TestCheckRuffSkip(t *testing.T) {
	if _, err := exec.LookPath("ruff"); err == nil {
		t.Skip("ruff is installed; skip not-installed test")
	}
	violations := CheckRuff([]string{"testdata/py_ruff_violation.py"})
	if len(violations) != 0 {
		t.Errorf("expected no violations when ruff is not installed, got %d", len(violations))
	}
}
