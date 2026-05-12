//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckPrettierRun — prettier 설치 시 포매팅 위반 파일에 N4 violation 생성 확인
package validate

import (
	"os/exec"
	"testing"
)

func TestCheckPrettierRun(t *testing.T) {
	if _, err := exec.LookPath("prettier"); err != nil {
		t.Skip("prettier not installed; skipping")
	}
	violations := CheckPrettier([]string{"testdata/ts_unformatted.ts"})
	if len(violations) == 0 {
		t.Error("expected N4 violation for unformatted file, got none")
		return
	}
	if violations[0].Rule != "N4" {
		t.Errorf("expected rule N4, got %s", violations[0].Rule)
	}
}
