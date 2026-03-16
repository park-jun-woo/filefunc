//ff:func feature=cli type=util control=iteration
//ff:what 대상 경로에서 상위로 탐색하며 go.mod 파일 경로를 찾음
//ff:checked llm=gpt-oss:20b hash=65fbd539
package cli

import (
	"os"
	"path/filepath"
)

// FindGoMod searches upward from target to find the nearest go.mod file.
func FindGoMod(target string) string {
	abs, err := filepath.Abs(target)
	if err != nil {
		return "go.mod"
	}
	for {
		candidate := filepath.Join(abs, "go.mod")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			break
		}
		abs = parent
	}
	return "go.mod"
}
