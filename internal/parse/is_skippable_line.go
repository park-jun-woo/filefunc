//ff:func feature=parse type=util
//ff:what 어노테이션 파싱 시 건너뛸 수 있는 라인인지 판별
package parse

import "strings"

// IsSkippableLine returns true if the line should be skipped during annotation parsing
// (empty lines, regular comments, package declarations).
func IsSkippableLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "package ")
}
