//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckBlackClean — black 설치 시 포매팅된 파일에 위반 없음 확인
package validate

import (
	"os/exec"
	"testing"
)

func TestCheckBlackClean(t *testing.T) {
	if _, err := exec.LookPath("black"); err != nil {
		t.Skip("black not installed; skipping")
	}
	violations := CheckBlack([]string{"testdata/py_formatted.py"})
	if len(violations) != 0 {
		t.Errorf("expected no violations for formatted file, got %d", len(violations))
	}
}
