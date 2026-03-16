//ff:func feature=parse type=parser control=iteration
//ff:what "key1=val1 key2=val2" 형식 문자열을 맵으로 파싱
//ff:checked llm=gpt-oss:20b hash=e3c9290a
package parse

import "strings"

// ParseFuncPairs parses "key1=val1 key2=val2" into a map.
func ParseFuncPairs(value string) map[string]string {
	result := make(map[string]string)
	for _, pair := range strings.Fields(value) {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
