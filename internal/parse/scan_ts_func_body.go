//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what function 선언 이후의 TypeScript func body 줄을 다음 top-level 선언까지 수집
package parse

import (
	"bufio"
	"strings"
)

// scanTSFuncBody collects lines after the first function until the next top-level declaration or EOF.
func scanTSFuncBody(scanner *bufio.Scanner, lines []string) []string {
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if isNextTSTopLevel(trimmed) {
			break
		}
		lines = append(lines, line)
	}
	return lines
}
