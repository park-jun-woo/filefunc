//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckPrettierClean — prettier 설치 시 포매팅된 파일에 위반 없음 확인
package validate

import (
	"os/exec"
	"testing"
)

func TestCheckPrettierClean(t *testing.T) {
	if _, err := exec.LookPath("prettier"); err != nil {
		t.Skip("prettier not installed; skipping")
	}
	violations := CheckPrettier([]string{"testdata/ts_formatted.ts"})
	if len(violations) != 0 {
		t.Errorf("expected no violations for formatted file, got %d", len(violations))
	}
}
