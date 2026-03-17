//ff:func feature=chain type=parser control=iteration dimension=1
//ff:what package별 chain 결과 필터링 — qualified name에서 패키지를 추출하여 선별
package chain

// FilterByPackage keeps only results whose qualified name has the given package prefix.
func FilterByPackage(results []ChonResult, pkg string) []ChonResult {
	var filtered []ChonResult
	for _, r := range results {
		if pkgFromQualified(r.Name) == pkg {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
