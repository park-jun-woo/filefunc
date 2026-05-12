//ff:func feature=validate type=util control=selection
//ff:what SourceFile의 필드명으로 해당 항목 수를 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

func countField(sf model.SourceFile, field string) int {
	switch field {
	case "Funcs":
		return len(sf.GetFuncs())
	case "Types":
		return len(sf.GetTypes())
	case "Methods":
		return len(sf.GetMethods())
	}
	return 0
}
