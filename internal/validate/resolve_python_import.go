//ff:func feature=validate type=util control=iteration dimension=1
//ff:what 상대 import 모듈 문자열을 프로젝트 내 .py 파일 경로로 변환
package validate

import (
	"os"
	"path/filepath"
	"strings"
)

// resolvePythonImport resolves a Python import module string to a file path.
// Relative imports (starting with dots) are resolved relative to fromFile.
// Absolute imports (no dots) are treated as external and return "".
// Returns "" if the resolved file does not exist under root.
func resolvePythonImport(fromFile string, importModule string, root string) string {
	if importModule == "" {
		return ""
	}

	dots := 0
	for _, ch := range importModule {
		if ch == '.' {
			dots++
		} else {
			break
		}
	}

	if dots == 0 {
		return ""
	}

	rest := importModule[dots:]

	dir := filepath.Dir(fromFile)
	for i := 1; i < dots; i++ {
		dir = filepath.Dir(dir)
	}

	if rest == "" {
		return ""
	}

	parts := strings.Split(rest, ".")
	candidate := filepath.Join(dir, filepath.Join(parts...)) + ".py"

	absCandidate, err := filepath.Abs(candidate)
	if err != nil {
		return ""
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return ""
	}

	if !strings.HasPrefix(absCandidate, absRoot) {
		return ""
	}

	if _, err := os.Stat(absCandidate); err != nil {
		return ""
	}

	return absCandidate
}
