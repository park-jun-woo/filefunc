//ff:func feature=validate type=util control=iteration dimension=1
//ff:what 한 파일의 모듈 import 목록을 resolve하여 프로젝트 내 파일 경로 슬라이스 반환
package validate

// resolveFileImports resolves a file's module imports to project file paths.
// Unresolved (external) imports are excluded from the result.
func resolveFileImports(absPath string, moduleImports []string, root string) []string {
	var edges []string
	for _, mod := range moduleImports {
		resolved := resolvePythonImport(absPath, mod, root)
		if resolved != "" {
			edges = append(edges, resolved)
		}
	}
	return edges
}
