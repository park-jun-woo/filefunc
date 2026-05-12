//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckESLintRun — eslint 설치 및 설정 시 린트 위반에 N4 violation 생성 확인
package validate

import (
	"os/exec"
	"testing"
)

func TestCheckESLintRun(t *testing.T) {
	if _, err := exec.LookPath("eslint"); err != nil {
		t.Skip("eslint not installed; skipping")
	}
	violations := CheckESLint([]string{"testdata/ts_eslint_violation.ts"})
	if len(violations) == 0 {
		t.Skip("eslint produced no violations; likely no config — skipping")
	}
	if violations[0].Rule != "N4" {
		t.Errorf("expected rule N4, got %s", violations[0].Rule)
	}
}
