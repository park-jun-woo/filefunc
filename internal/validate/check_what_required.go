//ff:func feature=validate type=rule
//ff:what A3: func 또는 type 파일에 //ff:what 필수 검증
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckWhatRequired checks A3: files with funcs or types must have //ff:what annotation.
func CheckWhatRequired(gf *model.GoFile) []model.Violation {
	if gf.IsTest {
		return nil
	}
	if len(gf.Funcs) == 0 && len(gf.Types) == 0 {
		return nil
	}
	if gf.Annotation == nil || gf.Annotation.What == "" {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A3",
			Level:   "ERROR",
			Message: "file with func or type must have //ff:what annotation",
		}}
	}
	return nil
}
