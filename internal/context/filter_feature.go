//ff:func feature=context type=util control=iteration dimension=1
//ff:what GoFile 목록에서 지정된 feature 값에 해당하는 파일만 필터링
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// FilterFeature keeps GoFiles whose annotation feature matches any of the given features.
func FilterFeature(files []*model.GoFile, features []string) []*model.GoFile {
	featureSet := make(map[string]bool, len(features))
	for _, f := range features {
		featureSet[f] = true
	}
	var kept []*model.GoFile
	for _, gf := range files {
		if gf.Annotation == nil {
			continue
		}
		f := gf.Annotation.Func["feature"]
		if f == "" {
			f = gf.Annotation.Type["feature"]
		}
		if featureSet[f] {
			kept = append(kept, gf)
		}
	}
	return kept
}
