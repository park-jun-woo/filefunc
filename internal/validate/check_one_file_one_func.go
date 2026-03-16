//ff:func feature=validate type=rule control=sequence
//ff:what F1: 파일당 func 1개 검증
//ff:checked llm=gpt-oss:20b hash=f8640445
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckOneFileOneFunc checks F1: each file must contain at most one func.
// Exceptions: _test.go files (F5), const-only files (F7).
func CheckOneFileOneFunc(gf *model.GoFile) []model.Violation {
	if gf.IsTest {
		return nil
	}
	if IsConstOnly(gf) {
		return nil
	}
	if len(gf.Funcs) > 1 {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "F1",
			Level:   "ERROR",
			Message: "file contains multiple funcs; expected 1 file 1 func",
		}}
	}
	return nil
}
