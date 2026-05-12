//ff:func feature=context type=util control=iteration dimension=1
//ff:what SourceFile 어노테이션에서 key=value AND 매칭 필터
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// FilterSearch keeps SourceFiles whose annotation matches all key=value pairs (AND).
func FilterSearch(files []model.SourceFile, query map[string]string) []model.SourceFile {
	var kept []model.SourceFile
	for _, sf := range files {
		ann := sf.GetAnnotation()
		if ann == nil {
			continue
		}
		if matchAnnotation(ann, query) {
			kept = append(kept, sf)
		}
	}
	return kept
}
