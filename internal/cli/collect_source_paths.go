//ff:func feature=cli type=util control=iteration dimension=1
//ff:what SourceFile 슬라이스에서 파일 경로 문자열 슬라이스를 추출
package cli

import "github.com/park-jun-woo/filefunc/internal/model"

// collectSourcePaths extracts file paths from a slice of SourceFile.
func collectSourcePaths(files []model.SourceFile) []string {
	paths := make([]string, len(files))
	for i, f := range files {
		paths[i] = f.GetPath()
	}
	return paths
}
