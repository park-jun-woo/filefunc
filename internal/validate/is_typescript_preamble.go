//ff:func feature=validate type=util control=sequence
//ff:what TypeScript 파일에서 어노테이션 이전에 허용되는 preamble 줄인지 판별
package validate

import "strings"

// isTypeScriptPreamble returns true if the line is a blank line or comment
// (including //ff: annotations) that may appear before code starts.
func isTypeScriptPreamble(line string) bool {
	return line == "" ||
		strings.HasPrefix(line, "//") ||
		strings.HasPrefix(line, "/*") ||
		strings.HasPrefix(line, "*") ||
		strings.HasPrefix(line, "*/")
}
