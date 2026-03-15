//ff:func feature=parse type=util
//ff:what 단일 .ffignore 패턴과 경로를 매칭
package walk

import (
	"path/filepath"
	"strings"
)

func matchPattern(path string, name string, isDir bool, pattern string) bool {
	if strings.HasSuffix(pattern, "/") {
		return isDir && (name == strings.TrimSuffix(pattern, "/") || strings.Contains(path, pattern))
	}
	matched, _ := filepath.Match(pattern, name)
	return matched
}
