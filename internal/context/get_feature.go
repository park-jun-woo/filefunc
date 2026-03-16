//ff:func feature=context type=util control=sequence
//ff:what 함수명에 해당하는 GoFile의 feature 값을 반환
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// getFeature returns the feature value for a function name.
func getFeature(name string, fileMap map[string]*model.GoFile) string {
	gf := fileMap[name]
	if gf == nil || gf.Annotation == nil {
		return ""
	}
	f := gf.Annotation.Func["feature"]
	if f == "" {
		f = gf.Annotation.Type["feature"]
	}
	return f
}
