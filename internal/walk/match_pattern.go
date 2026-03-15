//ff:func feature=parse type=util
//ff:what 단일 .ffignore 패턴과 경로를 매칭
package walk

import (
	"path/filepath"
	"strings"
)

func matchPattern(path string, name string, isDir bool, pattern string) bool {
	if strings.HasSuffix(pattern, "/") {
		dirPattern := strings.TrimSuffix(pattern, "/")
		if !isDir {
			return false
		}
		if name == dirPattern {
			return true
		}
		return strings.HasSuffix(path, dirPattern) || strings.Contains(path+"/", pattern)
	}
	matched, _ := filepath.Match(pattern, name)
	return matched
}
