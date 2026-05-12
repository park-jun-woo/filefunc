//ff:func feature=cli type=util control=selection
//ff:what 언어에 따라 프로젝트 루트 유효성 검증 — Go는 go.mod 필수, Python은 codebook.yaml만 필수
package cli

import "fmt"

// CheckProjectRootForLang verifies the path is a valid project root for the given language.
func CheckProjectRootForLang(root string, lang string) error {
	switch lang {
	case "go":
		return CheckProjectRoot(root)
	case "python":
		return CheckProjectRootPython(root)
	default:
		return fmt.Errorf("unsupported language: %s", lang)
	}
}
