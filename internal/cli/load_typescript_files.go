//ff:func feature=cli type=loader control=iteration dimension=1
//ff:what TypeScript 파일(.ts, .tsx)을 순회하고 파싱하여 SourceFile 슬라이스로 반환
package cli

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
	"github.com/park-jun-woo/filefunc/internal/walk"
)

// LoadTypeScriptFiles walks and parses TypeScript files under root, returning them as SourceFile slice.
func LoadTypeScriptFiles(root string, ignorePatterns []string) ([]model.SourceFile, error) {
	tsPaths, err := walk.WalkFiles(root, ".ts", ignorePatterns)
	if err != nil {
		return nil, fmt.Errorf("walking typescript .ts files: %w", err)
	}

	tsxPaths, err := walk.WalkFiles(root, ".tsx", ignorePatterns)
	if err != nil {
		return nil, fmt.Errorf("walking typescript .tsx files: %w", err)
	}

	allPaths := append(tsPaths, tsxPaths...)
	tsFiles, err := parse.ParseTypeScriptFiles(allPaths, root)
	if err != nil {
		return nil, fmt.Errorf("parsing typescript files: %w", err)
	}

	var files []model.SourceFile
	for _, tf := range tsFiles {
		ann, _ := parse.ParseTypeScriptAnnotation(tf.Path)
		tf.Annotation = ann
		files = append(files, tf)
	}
	return files, nil
}
