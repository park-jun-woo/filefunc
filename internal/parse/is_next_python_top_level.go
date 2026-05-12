//ff:func feature=parse type=util control=sequence
//ff:what Python 소스의 다음 top-level 정의(def/class) 시작 여부를 판단
package parse

import "strings"

// isNextPythonTopLevel returns true if the trimmed line starts a new top-level def or class.
func isNextPythonTopLevel(trimmed string) bool {
	if trimmed == "" {
		return false
	}
	return strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "class ")
}
