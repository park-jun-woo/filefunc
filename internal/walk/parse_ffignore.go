//ff:func feature=parse type=loader
//ff:what .ffignore 파일을 읽어 패턴 목록을 반환 (없으면 빈 목록)
package walk

import (
	"bufio"
	"os"
	"strings"
)

// ParseFFIgnore reads a .ffignore file and returns the list of patterns.
// Returns empty slice if file doesn't exist.
func ParseFFIgnore(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var patterns []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns
}
