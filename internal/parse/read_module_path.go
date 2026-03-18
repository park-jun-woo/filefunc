//ff:func feature=parse type=loader control=iteration dimension=1
//ff:what go.mod에서 모듈 경로를 추출
package parse

import (
	"bufio"
	"fmt"
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
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("module directive not found in %s", goModPath)
}
