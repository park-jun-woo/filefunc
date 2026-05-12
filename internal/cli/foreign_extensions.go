//ff:func feature=validate type=util control=selection
//ff:what 현재 언어에 대해 타 언어 확장자 목록 반환
package cli

// foreignExtensions returns the file extensions that are foreign to the given language.
func foreignExtensions(lang string) []string {
	switch lang {
	case "go":
		return []string{".py", ".ts", ".tsx"}
	case "python":
		return []string{".go", ".ts", ".tsx"}
	case "typescript":
		return []string{".go", ".py"}
	default:
		return nil
	}
}
