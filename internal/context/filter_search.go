//ff:func feature=context type=util control=iteration dimension=1
//ff:what GoFile 어노테이션에서 key=value AND 매칭 필터
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// FilterSearch keeps GoFiles whose annotation matches all key=value pairs (AND).
func FilterSearch(files []*model.GoFile, query map[string]string) []*model.GoFile {
	var kept []*model.GoFile
	for _, gf := range files {
		if gf.Annotation == nil {
			continue
		}
		if matchAnnotation(gf.Annotation, query) {
			kept = append(kept, gf)
		}
	}
	return kept
}
