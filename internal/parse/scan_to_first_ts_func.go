//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what TypeScript 소스를 스캔하며 첫 번째 function 선언 줄을 찾아 반환
package parse

import (
	"bufio"
	"strings"
)

// scanToFirstTSFunc scans until it finds a line starting with a TypeScript function declaration.
// Recognizes: "export function", "function", "async function", "export async function",
// "export default function".
// Returns the full line text, or "" if not found.
func scanToFirstTSFunc(scanner *bufio.Scanner) string {
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if isTSFuncStart(trimmed) {
			return line
		}
	}
	return ""
}
