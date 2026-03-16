//ff:func feature=context type=util control=iteration dimension=1
//ff:what chon=2 결과에서 대상 함수와 같은 feature만 유지, chon=1은 무조건 유지
package context

import (
	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/park-jun-woo/filefunc/internal/model"
)

// FilterFeature keeps chon=1 unconditionally and filters chon=2 to same feature.
func FilterFeature(results []chain.ChonResult, targetFeature string, fileMap map[string]*model.GoFile) []chain.ChonResult {
	var kept []chain.ChonResult
	for _, r := range results {
		if r.Chon <= 1 {
			kept = append(kept, r)
			continue
		}
		gf := fileMap[r.Name]
		if gf == nil || gf.Annotation == nil {
			continue
		}
		f := gf.Annotation.Func["feature"]
		if f == "" {
			f = gf.Annotation.Type["feature"]
		}
		if f == targetFeature {
			kept = append(kept, r)
		}
	}
	return kept
}
