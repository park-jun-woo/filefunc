//ff:func feature=parse type=walker control=sequence
//ff:what 프로젝트 루트의 설정 파일로 프로그래밍 언어를 자동 감지
package walk

// DetectLang detects the project language from marker files in the root directory.
// Returns "go", "python", or "" (empty string if undetermined).
func DetectLang(root string) string {
	if fileExists(root, "go.mod") {
		return "go"
	}
	if fileExists(root, "setup.py") || fileExists(root, "pyproject.toml") || fileExists(root, "setup.cfg") {
		return "python"
	}
	return ""
}
