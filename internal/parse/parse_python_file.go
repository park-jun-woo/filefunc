//ff:func feature=parse type=parser control=sequence
//ff:what 단일 Python 파일을 파싱하여 PythonFile 반환
package parse

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// ParsePythonFile parses a single Python file and returns a PythonFile.
func ParsePythonFile(path string) (*model.PythonFile, error) {
	files, err := ParsePythonFiles([]string{path})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no result for %s", path)
	}
	return files[0], nil
}
