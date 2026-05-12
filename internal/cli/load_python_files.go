//ff:func feature=cli type=loader control=iteration dimension=1
//ff:what Python 파일을 순회하고 파싱하여 SourceFile 슬라이스로 반환
package cli

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
	"github.com/park-jun-woo/filefunc/internal/walk"
)

// LoadPythonFiles walks and parses Python files under root, returning them as SourceFile slice.
func LoadPythonFiles(root string, ignorePatterns []string) ([]model.SourceFile, error) {
	pyPaths, err := walk.WalkFiles(root, ".py", ignorePatterns)
	if err != nil {
		return nil, fmt.Errorf("walking python files: %w", err)
	}

	pyFiles, err := parse.ParsePythonFiles(pyPaths)
	if err != nil {
		return nil, fmt.Errorf("parsing python files: %w", err)
	}

	var files []model.SourceFile
	for _, pf := range pyFiles {
		ann, _ := parse.ParsePythonAnnotation(pf.Path)
		pf.Annotation = ann
		files = append(files, pf)
	}
	return files, nil
}
