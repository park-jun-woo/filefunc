//ff:func feature=validate type=rule control=sequence
//ff:what F3: 파일당 method 1개 검증
//ff:checked llm=gpt-oss:20b hash=2e1f09f6
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckOneFileOneMethod checks F3: each file must contain at most one method.
func CheckOneFileOneMethod(gf *model.GoFile) []model.Violation {
	if gf.IsTest {
		return nil
	}
	if len(gf.Methods) > 1 {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "F3",
			Level:   "ERROR",
			Message: "file contains multiple methods; expected 1 file 1 method",
		}}
	}
	return nil
}
