//ff:func feature=parse type=util control=selection
//ff:what ts_ast.js 실행 에러를 분석하여 사용자 친화적 메시지로 변환
package parse

import (
	"fmt"
	"os/exec"
	"strings"
)

func wrapTsAstError(err error) error {
	ee, ok := err.(*exec.ExitError)
	if !ok {
		return fmt.Errorf("ts_ast.js exec: %w", err)
	}

	stderr := string(ee.Stderr)
	switch {
	case strings.Contains(stderr, "MODULE_NOT_FOUND"),
		strings.Contains(stderr, "Cannot find module"):
		return fmt.Errorf("typescript package not found in node_modules; run: npm install typescript")
	default:
		return fmt.Errorf("ts_ast.js failed: %s", stderr)
	}
}
