//ff:func feature=parse type=parser control=sequence
//ff:what Python 파일에서 첫 번째 func/method의 전체 소스를 텍스트 기반으로 추출
package parse

import (
	"bufio"
	"os"
	"strings"
)

// ExtractFuncSourcePython extracts the full source of the first def/method
// from a Python file. Returns empty string if no def found.
func ExtractFuncSourcePython(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	defLine := scanToFirstDef(scanner)
	if defLine == "" {
		return "", scanner.Err()
	}

	lines := []string{defLine}
	lines = scanFuncBody(scanner, lines)
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}
