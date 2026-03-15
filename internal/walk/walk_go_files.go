//ff:func feature=parse type=walker
//ff:what 디렉토리를 재귀 순회하며 .go 파일 경로 목록 반환 (.ffignore 적용)
package walk

import (
	"os"
	"path/filepath"
	"strings"
)

// WalkGoFiles recursively walks a directory and returns all .go file paths.
// Skips testdata directories and paths matching ignorePatterns.
func WalkGoFiles(root string, ignorePatterns []string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "testdata" {
			return filepath.SkipDir
		}
		if MatchFFIgnore(path, info.Name(), info.IsDir(), ignorePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
