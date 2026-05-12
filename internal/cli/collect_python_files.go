//ff:func feature=cli type=util control=iteration dimension=1
//ff:what SourceFile 슬라이스에서 PythonFile 포인터 슬라이스를 추출
package cli

import "github.com/park-jun-woo/filefunc/internal/model"

// collectPythonFiles extracts *PythonFile from a slice of SourceFile.
func collectPythonFiles(files []model.SourceFile) []*model.PythonFile {
	var result []*model.PythonFile
	for _, sf := range files {
		if pf, ok := sf.(*model.PythonFile); ok {
			result = append(result, pf)
		}
	}
	return result
}
