//ff:func feature=cli type=util control=selection
//ff:what 언어에 따라 func body 해시를 계산 — Go는 AST 기반, Python은 PythonFile.BodyHash 사용
package cli

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// CalcBodyHashForLang returns the body hash for the source file, dispatching by language.
func CalcBodyHashForLang(sf model.SourceFile) (string, error) {
	switch sf.GetLang() {
	case "python":
		pf, ok := sf.(*model.PythonFile)
		if !ok {
			return "", nil
		}
		return pf.BodyHash, nil
	default:
		return parse.CalcBodyHash(sf.GetPath())
	}
}
