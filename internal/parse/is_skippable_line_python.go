//ff:func feature=parse type=util control=sequence
//ff:what Python 어노테이션 파싱 시 건너뛸 수 있는 라인인지 판별
package parse

import "strings"

// IsSkippableLinePython returns true if the line should be skipped during Python annotation parsing.
// Skips empty lines, shebang, encoding declarations, and non-ff comments.
func IsSkippableLinePython(line string) bool {
	if line == "" {
		return true
	}
	if strings.HasPrefix(line, "#!/") {
		return true
	}
	if strings.HasPrefix(line, "# coding:") || strings.HasPrefix(line, "# -*- coding") {
		return true
	}
	if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "# ff:") {
		return true
	}
	return false
}
