//ff:func feature=annotate type=util
//ff:what 어노테이션 블록 뒤에 새 라인을 삽입
//ff:checked llm=gpt-oss:20b hash=62db7a3d
package annotate

import "strings"

// InsertAfterAnnotations inserts a new line after the last //ff: annotation line.
func InsertAfterAnnotations(lines []string, newLine string) []string {
	insertIdx := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//ff:") {
			insertIdx = i + 1
		}
	}
	var result []string
	result = append(result, lines[:insertIdx]...)
	result = append(result, newLine)
	result = append(result, lines[insertIdx:]...)
	return result
}
