//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckESLintNotInstalled — eslint 미설치 시 N4 ERROR 반환 확인
package validate

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCheckESLintNotInstalled(t *testing.T) {
	if _, err := exec.LookPath("eslint"); err == nil {
		t.Skip("eslint is installed; skip not-installed test")
	}
	violations := CheckESLint([]string{"testdata/ts_eslint_violation.ts"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "eslint not installed") {
		t.Errorf("expected 'eslint not installed' message, got %q", violations[0].Message)
	}
}
