//ff:func feature=validate type=util control=sequence
//ff:what test: TestCheckPrettierNotInstalled — prettier 미설치 시 N4 ERROR 반환 확인
package validate

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCheckPrettierNotInstalled(t *testing.T) {
	if _, err := exec.LookPath("prettier"); err == nil {
		t.Skip("prettier is installed; skip not-installed test")
	}
	violations := CheckPrettier([]string{"testdata/ts_unformatted.ts"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !strings.Contains(violations[0].Message, "prettier not installed") {
		t.Errorf("expected 'prettier not installed' message, got %q", violations[0].Message)
	}
}
