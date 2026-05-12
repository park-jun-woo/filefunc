//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckBlackRun — black 설치 시 포매팅 위반 파일에 N4 violation 생성 확인
package validate

import (
	"os/exec"
	"testing"
)

func TestCheckBlackRun(t *testing.T) {
	if _, err := exec.LookPath("black"); err != nil {
		t.Skip("black not installed; skipping")
	}
	violations := CheckBlack([]string{"testdata/py_unformatted.py"})
	if len(violations) == 0 {
		t.Error("expected N4 violation for unformatted file, got none")
		return
	}
	if violations[0].Rule != "N4" {
		t.Errorf("expected rule N4, got %s", violations[0].Rule)
	}
}
