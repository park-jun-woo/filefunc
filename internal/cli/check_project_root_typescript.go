//ff:func feature=cli type=util control=sequence
//ff:what TypeScript 프로젝트 루트 유효성 검증 — 디렉토리, tsconfig.json, codebook.yaml 존재 확인
package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckProjectRootTypeScript verifies the path is a directory with tsconfig.json and codebook.yaml.
func CheckProjectRootTypeScript(root string) error {
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("path not found: %s", root)
	}
	if !info.IsDir() {
		return fmt.Errorf("path must be a directory: %s", root)
	}
	if _, err := os.Stat(filepath.Join(root, "tsconfig.json")); err != nil {
		return fmt.Errorf("tsconfig.json not found in %s", root)
	}
	if _, err := os.Stat(filepath.Join(root, "codebook.yaml")); err != nil {
		return fmt.Errorf("codebook.yaml not found in %s", root)
	}
	return nil
}
