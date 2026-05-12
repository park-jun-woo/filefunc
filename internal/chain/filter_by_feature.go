//ff:func feature=chain type=parser control=iteration dimension=1
//ff:what feature별 func 필터링 — 어노테이션의 feature 값으로 선별
package chain

import "github.com/park-jun-woo/filefunc/internal/model"

// FilterByFeature returns func names whose //ff:func or //ff:type annotation control=iteration dimension=1
// has the given feature value.
func FilterByFeature(files []model.SourceFile, feature string) []string {
	var result []string
	for _, sf := range files {
		if !hasFeature(sf, feature) {
			continue
		}
		for _, name := range sf.GetFuncs() {
			result = append(result, qualifiedName(sf.GetPackage(), name))
		}
	}
	return result
}
