//ff:func feature=validate type=util control=selection
//ff:what GoFile의 필드명으로 해당 항목 수를 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

func countField(gf *model.GoFile, field string) int {
	switch field {
	case "Funcs":
		return len(gf.Funcs)
	case "Types":
		return len(gf.Types)
	case "Methods":
		return len(gf.Methods)
	}
	return 0
}
