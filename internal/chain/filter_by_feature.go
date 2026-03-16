//ff:func feature=chain type=parser control=iteration
//ff:what feature별 func 필터링 — 어노테이션의 feature 값으로 선별
package chain

import "github.com/park-jun-woo/filefunc/internal/model"

// FilterByFeature returns func names whose //ff:func or //ff:type annotation control=iteration
// has the given feature value.
func FilterByFeature(files []*model.GoFile, feature string) []string {
	var result []string
	for _, gf := range files {
		if !hasFeature(gf, feature) {
			continue
		}
		result = append(result, gf.Funcs...)
	}
	return result
}
