//ff:func feature=validate type=util control=selection
//ff:what Python 파일에서 어노테이션 이전에 허용되는 preamble 줄인지 판별
package validate

import "strings"

// isPythonPreamble returns true if the line is a shebang, encoding declaration,
// blank line, or comment (including # ff: annotations).
func isPythonPreamble(line string) bool {
	switch {
	case line == "":
		return true
	case strings.HasPrefix(line, "#!"):
		return true
	case strings.HasPrefix(line, "# coding:"), strings.HasPrefix(line, "# -*- coding"):
		return true
	case strings.HasPrefix(line, "#"):
		return true
	}
	return false
}
