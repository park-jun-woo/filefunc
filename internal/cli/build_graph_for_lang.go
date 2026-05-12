//ff:func feature=cli type=util control=selection
//ff:what 언어에 따라 Go, Python, TypeScript 호출 그래프 구성을 디스패치
package cli

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/park-jun-woo/filefunc/internal/model"
)

// BuildGraphForLang dispatches graph building based on language.
func BuildGraphForLang(root string, lang string) (*chain.CallGraph, []model.SourceFile, error) {
	switch lang {
	case "go":
		return BuildGraph(root)
	case "python":
		return BuildGraphPython(root)
	case "typescript":
		return BuildGraphTypeScript(root)
	default:
		return nil, nil, fmt.Errorf("unsupported language: %s", lang)
	}
}
