//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckESLintClean — eslint 설치 및 설정 시 위반 없는 파일에 빈 결과 확인
package validate

import (
	"os/exec"
	"testing"
)

func TestCheckESLintClean(t *testing.T) {
	if _, err := exec.LookPath("eslint"); err != nil {
		t.Skip("eslint not installed; skipping")
	}
	violations := CheckESLint([]string{"testdata/ts_formatted.ts"})
	if len(violations) != 0 {
		t.Errorf("expected no violations for clean file, got %d", len(violations))
	}
}
