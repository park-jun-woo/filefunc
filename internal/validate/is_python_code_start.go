//ff:func feature=validate type=util control=sequence
//ff:what Python 코드 시작 줄인지 판별 (import, def, class, from 등)
package validate

import "strings"

// isPythonCodeStart returns true if the line starts Python code.
func isPythonCodeStart(line string) bool {
	return strings.HasPrefix(line, "import ") ||
		strings.HasPrefix(line, "from ") ||
		strings.HasPrefix(line, "def ") ||
		strings.HasPrefix(line, "class ") ||
		strings.HasPrefix(line, "async def ") ||
		strings.HasPrefix(line, "async class ")
}
