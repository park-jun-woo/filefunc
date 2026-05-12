//ff:func feature=validate type=util control=sequence
//ff:what TypeScript 코드 시작 줄인지 판별 (import, export, class, function, const, let, var, enum, interface, type, async 등)
package validate

import "strings"

// isTypeScriptCodeStart returns true if the line starts TypeScript code.
func isTypeScriptCodeStart(line string) bool {
	return strings.HasPrefix(line, "import ") ||
		strings.HasPrefix(line, "export ") ||
		strings.HasPrefix(line, "class ") ||
		strings.HasPrefix(line, "function ") ||
		strings.HasPrefix(line, "const ") ||
		strings.HasPrefix(line, "let ") ||
		strings.HasPrefix(line, "var ") ||
		strings.HasPrefix(line, "enum ") ||
		strings.HasPrefix(line, "interface ") ||
		strings.HasPrefix(line, "type ") ||
		strings.HasPrefix(line, "async ")
}
