//ff:func feature=annotate type=formatter control=sequence
//ff:what 기존 //ff: 라인을 새 값으로 교체하거나 제거하여 결과 라인 반환
//ff:checked llm=gpt-oss:20b hash=20f78335
package annotate

import "strings"

// ReplaceAnnotationLine processes a single line, replacing or removing matching //ff: annotations.
// Returns the replacement line (empty string means remove), and whether a match was found.
func ReplaceAnnotationLine(trimmed string, prefix string, key string, newLine string, value string) (replacement string, matched bool) {
	if !strings.HasPrefix(trimmed, prefix) && trimmed != "//ff:"+key {
		return trimmed, false
	}
	if value == "" {
		return "", true
	}
	return newLine, true
}
