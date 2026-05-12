//ff:func feature=parse type=parser control=sequence
//ff:what 단일 TypeScript 파일을 파싱하여 TypeScriptFile 반환
package parse

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// ParseTypeScriptFile parses a single TypeScript file and returns a TypeScriptFile.
func ParseTypeScriptFile(path string, projectRoot string) (*model.TypeScriptFile, error) {
	files, err := ParseTypeScriptFiles([]string{path}, projectRoot)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no result for %s", path)
	}
	return files[0], nil
}
