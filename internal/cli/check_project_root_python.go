//ff:func feature=cli type=util control=sequence
//ff:what Python 프로젝트 루트 유효성 검증 — 디렉토리와 codebook.yaml 존재 확인
package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckProjectRootPython verifies the path is a directory with codebook.yaml.
func CheckProjectRootPython(root string) error {
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("path not found: %s", root)
	}
	if !info.IsDir() {
		return fmt.Errorf("path must be a directory: %s", root)
	}
	if _, err := os.Stat(filepath.Join(root, "codebook.yaml")); err != nil {
		return fmt.Errorf("codebook.yaml not found in %s", root)
	}
	return nil
}
