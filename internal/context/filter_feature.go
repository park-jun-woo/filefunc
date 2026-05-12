//ff:func feature=context type=util control=iteration dimension=1
//ff:what SourceFile 목록에서 지정된 feature 값에 해당하는 파일만 필터링
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// FilterFeature keeps SourceFiles whose annotation feature matches any of the given features.
func FilterFeature(files []model.SourceFile, features []string) []model.SourceFile {
	featureSet := make(map[string]bool, len(features))
	for _, f := range features {
		featureSet[f] = true
	}
	var kept []model.SourceFile
	for _, sf := range files {
		ann := sf.GetAnnotation()
		if ann == nil {
			continue
		}
		f := ann.Func["feature"]
		if f == "" {
			f = ann.Type["feature"]
		}
		if featureSet[f] {
			kept = append(kept, sf)
		}
	}
	return kept
}
