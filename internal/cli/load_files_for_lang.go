//ff:func feature=cli type=loader control=selection
//ff:what 언어에 따라 Go 또는 Python 파일 로딩을 디스패치
package cli

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// LoadFilesForLang dispatches file loading based on language.
func LoadFilesForLang(root string, lang string, ignorePatterns []string) ([]model.SourceFile, error) {
	switch lang {
	case "go":
		return LoadGoFiles(root, ignorePatterns)
	case "python":
		return LoadPythonFiles(root, ignorePatterns)
	default:
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
}
