//ff:func feature=parse type=walker control=sequence
//ff:what 디렉토리를 재귀 순회하며 .go 파일 경로 목록 반환 (WalkFiles 래퍼)
package walk

// WalkGoFiles recursively walks a directory and returns all .go file paths.
// Delegates to WalkFiles with ".go" extension.
func WalkGoFiles(root string, ignorePatterns []string) ([]string, error) {
	return WalkFiles(root, ".go", ignorePatterns)
}
