//ff:func feature=cli type=util
//ff:what 지정 경로가 유효한 프로젝트 루트인지 검증 (폴더, go.mod, codebook.yaml 필수)
package cli

import (
	"fmt"
	"os"
)

// CheckProjectRoot verifies the path is a directory with go.mod and codebook.yaml.
func CheckProjectRoot(root string) error {
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("path not found: %s", root)
	}
	if !info.IsDir() {
		return fmt.Errorf("path must be a directory: %s", root)
	}
	if _, err := os.Stat(root + "/go.mod"); err != nil {
		return fmt.Errorf("go.mod not found in %s", root)
	}
	if _, err := os.Stat(root + "/codebook.yaml"); err != nil {
		return fmt.Errorf("codebook.yaml not found in %s", root)
	}
	return nil
}
