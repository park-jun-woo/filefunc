//ff:func feature=cli type=util control=sequence
//ff:what 대상 경로에서 상위로 탐색하며 go.mod가 있는 디렉토리 경로를 반환
package cli

import "path/filepath"

// FindGoModDir returns the directory containing go.mod by searching upward from target.
func FindGoModDir(target string) string {
	goModPath := FindGoMod(target)
	return filepath.Dir(goModPath)
}
