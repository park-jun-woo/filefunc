//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what Python 소스를 스캔하며 첫 번째 def 줄을 찾아 반환
package parse

import (
	"bufio"
	"strings"
)

// scanToFirstDef scans until it finds a line starting with "def ".
// Returns the full line text, or "" if not found.
func scanToFirstDef(scanner *bufio.Scanner) string {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "def ") {
			return line
		}
	}
	return ""
}
