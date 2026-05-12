//ff:func feature=chain type=parser control=selection
//ff:what 언어별(Go/Python/TypeScript) 호출 목록을 추출하는 디스패처
package chain

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// extractCallsForFile dispatches call extraction based on language.
func extractCallsForFile(sf model.SourceFile, modulePath string, projFuncs map[string]string) ([]string, error) {
	switch sf.GetLang() {
	case "go":
		return parse.ExtractCalls(sf.GetPath(), modulePath, projFuncs, sf.GetPackage())
	case "python":
		pf, ok := sf.(*model.PythonFile)
		if !ok {
			return nil, nil
		}
		return pf.Calls, nil
	case "typescript":
		tf, ok := sf.(*model.TypeScriptFile)
		if !ok {
			return nil, nil
		}
		return tf.Calls, nil
	default:
		return nil, nil
	}
}
