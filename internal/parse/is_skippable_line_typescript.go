//ff:func feature=parse type=util control=sequence
//ff:what TypeScript 어노테이션 파싱 시 건너뛸 수 있는 라인인지 판별
package parse

import "strings"

// IsSkippableLineTypeScript returns true if the line should be skipped during TypeScript annotation parsing.
// Skips empty lines and non-ff comments.
func IsSkippableLineTypeScript(line string) bool {
	if line == "" {
		return true
	}
	if strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "//ff:") {
		return true
	}
	return false
}
