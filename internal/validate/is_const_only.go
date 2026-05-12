//ff:func feature=validate type=util control=sequence
//ff:what 파일이 const/var만 포함하는지 판별
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// IsConstOnly returns true if the file contains only const/var declarations
// with no funcs, types, methods, or init.
func IsConstOnly(sf model.SourceFile) bool {
	return len(sf.GetFuncs()) == 0 && len(sf.GetTypes()) == 0 && len(sf.GetMethods()) == 0 && !sf.GetHasInit()
}
