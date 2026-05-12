//ff:func feature=cli type=util control=selection
//ff:what 언어에 따라 프로젝트 루트 유효성 검증 — Go는 go.mod, TypeScript는 tsconfig.json+codebook.yaml, Python은 codebook.yaml
package cli

import "fmt"

// CheckProjectRootForLang verifies the path is a valid project root for the given language.
func CheckProjectRootForLang(root string, lang string) error {
	switch lang {
	case "go":
		return CheckProjectRoot(root)
	case "python":
		return CheckProjectRootPython(root)
	case "typescript":
		return CheckProjectRootTypeScript(root)
	default:
		return fmt.Errorf("unsupported language: %s", lang)
	}
}
