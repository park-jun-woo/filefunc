//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what def 이후의 Python func body 줄을 다음 top-level 정의까지 수집
package parse

import (
	"bufio"
	"strings"
)

// scanFuncBody collects lines after the first def until the next top-level def/class/EOF.
func scanFuncBody(scanner *bufio.Scanner, lines []string) []string {
	for scanner.Scan() {
		line := scanner.Text()
		if isNextPythonTopLevel(strings.TrimSpace(line)) {
			break
		}
		lines = append(lines, line)
	}
	return lines
}
