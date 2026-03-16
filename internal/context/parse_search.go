//ff:func feature=context type=parser control=iteration dimension=1
//ff:what --search 문자열을 key=value 맵으로 파싱
package context

import "strings"

// ParseSearch parses "feature=X ssot=Y" into a map.
func ParseSearch(search string) map[string]string {
	result := make(map[string]string)
	for _, pair := range strings.Fields(search) {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
