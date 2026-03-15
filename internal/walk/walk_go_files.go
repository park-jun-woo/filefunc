//ff:func feature=parse type=walker
//ff:what 디렉토리를 재귀 순회하며 .go 파일 경로 목록 반환
//ff:checked llm=gpt-oss:20b hash=1791bca7
package walk

import (
	"os"
	"path/filepath"
	"strings"
)

// WalkGoFiles recursively walks a directory and returns all .go file paths.
// Skips testdata directories. Test files (_test.go) are included but can be
// identified by suffix.
func WalkGoFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "testdata" {
			return filepath.SkipDir
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
