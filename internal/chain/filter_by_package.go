//ff:func feature=chain type=parser control=iteration dimension=1
//ff:what package별 chain 결과 필터링 — 시작 노드 검증 + 결과에서 해당 패키지만 선별
package chain

import "github.com/park-jun-woo/filefunc/internal/model"

// FilterByPackage keeps only results whose GoFile.Package matches pkg.
func FilterByPackage(results []ChonResult, fileMap map[string]*model.GoFile, pkg string) []ChonResult {
	var filtered []ChonResult
	for _, r := range results {
		gf, ok := fileMap[r.Name]
		if ok && gf.Package == pkg {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
