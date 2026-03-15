//ff:func feature=validate type=util
//ff:what 파일이 const/var만 포함하는지 판별
//ff:checked llm=gpt-oss:20b hash=4deabf19
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// IsConstOnly returns true if the file contains only const/var declarations
// with no funcs, types, methods, or init.
func IsConstOnly(gf *model.GoFile) bool {
	return len(gf.Funcs) == 0 && len(gf.Types) == 0 && len(gf.Methods) == 0 && !gf.HasInit
}
