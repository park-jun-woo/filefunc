//ff:func feature=chain type=util control=iteration
//ff:what 그래프에서 관계 목록을 ChonResult로 수집하되 중복 제거
package chain

// CollectChon appends funcs from names to results if not already seen.
func CollectChon(names []string, chon int, rel string, seen map[string]bool) []ChonResult {
	var results []ChonResult
	for _, name := range names {
		if seen[name] {
			continue
		}
		results = append(results, ChonResult{name, chon, rel})
		seen[name] = true
	}
	return results
}
