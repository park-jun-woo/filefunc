//ff:func feature=parse type=util control=sequence
//ff:what trimmed line이 TypeScript function 선언 시작인지 판별
package parse

import "strings"

// isTSFuncStart returns true if the trimmed line starts a TypeScript function declaration.
func isTSFuncStart(trimmed string) bool {
	if strings.HasPrefix(trimmed, "function ") {
		return true
	}
	if strings.HasPrefix(trimmed, "async function ") {
		return true
	}
	if strings.HasPrefix(trimmed, "export function ") {
		return true
	}
	if strings.HasPrefix(trimmed, "export async function ") {
		return true
	}
	if strings.HasPrefix(trimmed, "export default function") {
		return true
	}
	return false
}
