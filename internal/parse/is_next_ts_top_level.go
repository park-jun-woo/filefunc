//ff:func feature=parse type=util control=sequence
//ff:what TypeScript 소스의 다음 top-level 선언(function/class/interface/type/export/const/let/var) 시작 여부를 판단
package parse

import "strings"

// isNextTSTopLevel returns true if the trimmed line starts a new top-level declaration.
func isNextTSTopLevel(trimmed string) bool {
	if trimmed == "" {
		return false
	}
	if isTSFuncStart(trimmed) {
		return true
	}
	if strings.HasPrefix(trimmed, "class ") || strings.HasPrefix(trimmed, "export class ") {
		return true
	}
	if strings.HasPrefix(trimmed, "interface ") || strings.HasPrefix(trimmed, "export interface ") {
		return true
	}
	if strings.HasPrefix(trimmed, "type ") || strings.HasPrefix(trimmed, "export type ") {
		return true
	}
	if strings.HasPrefix(trimmed, "export default ") {
		return true
	}
	return false
}
