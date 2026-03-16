//ff:func feature=parse type=loader control=iteration
//ff:what go.mod에서 모듈 경로를 추출
//ff:checked llm=gpt-oss:20b hash=2c02ae35
package parse

import (
	"bufio"
	"os"
	"strings"
)

// ReadModulePath reads the module path from a go.mod file.
func ReadModulePath(goModPath string) (string, error) {
	f, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(line[len("module "):]), nil
		}
	}
	return "", scanner.Err()
}
