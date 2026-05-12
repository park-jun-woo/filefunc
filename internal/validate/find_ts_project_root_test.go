//ff:func feature=validate type=util control=iteration dimension=1
//ff:what test: findTsProjectRoot — 테스트에서 node_modules가 있는 프로젝트 루트를 찾는 헬퍼
package validate

import (
	"os"
	"path/filepath"
)

func findTsProjectRoot() string {
	dir, _ := os.Getwd()
	for dir != "/" {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	return ""
}
