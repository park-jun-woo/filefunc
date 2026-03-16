//ff:func feature=chain type=util control=iteration
//ff:what 중복 없이 문자열을 결과 슬라이스에 추가
package chain

// AddUnique appends items to result if not already in seen.
func AddUnique(items []string, seen map[string]bool, result *[]string) {
	for _, item := range items {
		if seen[item] {
			continue
		}
		seen[item] = true
		*result = append(*result, item)
	}
}
