//ff:func feature=validate type=engine control=selection
//ff:what 소스 파일 언어에 따라 적절한 toulmin validate graph를 선택
package validate

import (
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// selectGraph returns the appropriate validate graph based on the file's language.
func selectGraph(sf model.SourceFile) *toulmin.Graph {
	switch sf.GetLang() {
	case "python":
		return PythonValidateGraph
	case "typescript":
		return TypeScriptValidateGraph
	default:
		return ValidateGraph
	}
}
