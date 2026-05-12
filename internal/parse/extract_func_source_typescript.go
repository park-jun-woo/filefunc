//ff:func feature=parse type=parser control=sequence
//ff:what TypeScript 파일에서 첫 번째 function/method body를 텍스트 기반으로 추출
package parse

import (
	"bufio"
	"os"
	"strings"
)

// ExtractFuncSourceTypeScript extracts the full source of the first function/method
// from a TypeScript file. Returns empty string if no function found.
func ExtractFuncSourceTypeScript(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	funcLine := scanToFirstTSFunc(scanner)
	if funcLine == "" {
		return "", scanner.Err()
	}

	lines := []string{funcLine}
	lines = scanTSFuncBody(scanner, lines)
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}
