//ff:func feature=parse type=util control=iteration dimension=1
//ff:what 쉼표 구분 문자열을 분리하고 각 항목 공백 제거
//ff:checked llm=gpt-oss:20b hash=addbd68c
package parse

import "strings"

// SplitTrim splits a comma-separated string and trims whitespace from each item.
func SplitTrim(s string) []string {
	var result []string
	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}
