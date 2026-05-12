//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckRuffClean — ruff 설치 시 위반 없는 파일에 빈 결과 확인
package validate

import (
	"os/exec"
	"testing"
)

func TestCheckRuffClean(t *testing.T) {
	if _, err := exec.LookPath("ruff"); err != nil {
		t.Skip("ruff not installed; skipping")
	}
	violations := CheckRuff([]string{"testdata/py_formatted.py"})
	if len(violations) != 0 {
		t.Errorf("expected no violations for clean file, got %d", len(violations))
	}
}
