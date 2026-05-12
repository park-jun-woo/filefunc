//ff:func feature=cli type=loader control=iteration dimension=1
//ff:what Go 파일을 순회하고 파싱하여 SourceFile 슬라이스로 반환
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
	"github.com/park-jun-woo/filefunc/internal/walk"
)

// LoadGoFiles walks and parses Go files under root, returning them as SourceFile slice.
func LoadGoFiles(root string, ignorePatterns []string) ([]model.SourceFile, error) {
	paths, err := walk.WalkGoFiles(root, ignorePatterns)
	if err != nil {
		return nil, fmt.Errorf("walking files: %w", err)
	}

	var files []model.SourceFile
	for _, p := range paths {
		gf, err := parse.ParseGoFile(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", p, err)
			continue
		}
		files = append(files, gf)
	}
	return files, nil
}
