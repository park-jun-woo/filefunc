//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckBlackNotInstalled — black 미설치 시 N4 ERROR 반환 확인
package validate

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCheckBlackNotInstalled(t *testing.T) {
	if _, err := exec.LookPath("black"); err == nil {
		t.Skip("black is installed; skip not-installed test")
	}
	violations := CheckBlack([]string{"testdata/py_unformatted.py"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "black not installed") {
		t.Errorf("expected 'black not installed' message, got %q", violations[0].Message)
	}
}
