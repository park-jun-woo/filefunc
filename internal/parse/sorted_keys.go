//ff:func feature=parse type=util
//ff:what bool 맵의 키를 정렬된 문자열 슬라이스로 반환
//ff:checked llm=gpt-oss:20b hash=3037cbab
package parse

import "sort"

// SortedKeys returns the keys of a map[string]bool as a sorted slice.
func SortedKeys(m map[string]bool) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
