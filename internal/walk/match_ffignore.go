//ff:func feature=parse type=util control=iteration
//ff:what 경로가 .ffignore 패턴에 매칭되는지 판별
package walk

// MatchFFIgnore returns true if the given path matches any of the ignore patterns.
func MatchFFIgnore(path string, name string, isDir bool, patterns []string) bool {
	for _, pattern := range patterns {
		if matchPattern(path, name, isDir, pattern) {
			return true
		}
	}
	return false
}
