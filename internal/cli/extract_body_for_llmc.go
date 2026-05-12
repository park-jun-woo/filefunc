//ff:func feature=cli type=util control=selection
//ff:what 언어에 따라 LLM 검증용 func body를 추출
package cli

import (
	"os"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// ExtractBodyForLlmc extracts the func body for LLM verification, dispatching by language.
func ExtractBodyForLlmc(sf model.SourceFile) string {
	switch sf.GetLang() {
	case "python":
		body, err := parse.ExtractFuncSourcePython(sf.GetPath())
		if err != nil {
			return ""
		}
		return body
	default:
		src, err := os.ReadFile(sf.GetPath())
		if err != nil {
			return ""
		}
		return parse.ExtractFuncSource(sf.GetPath(), src)
	}
}
