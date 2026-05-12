//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckBlackSkip — black 미설치 시 빈 슬라이스 반환 확인
package validate

import (
	"os/exec"
	"testing"
)

func TestCheckBlackSkip(t *testing.T) {
	if _, err := exec.LookPath("black"); err == nil {
		t.Skip("black is installed; skip not-installed test")
	}
	violations := CheckBlack([]string{"testdata/py_unformatted.py"})
	if len(violations) != 0 {
		t.Errorf("expected no violations when black is not installed, got %d", len(violations))
	}
}
