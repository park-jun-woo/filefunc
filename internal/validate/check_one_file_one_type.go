//ff:func feature=validate type=rule control=sequence
//ff:what F2: 파일당 type 1개 검증
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckOneFileOneType checks F2: each file must contain at most one type.
func CheckOneFileOneType(gf *model.GoFile) []model.Violation {
	if gf.IsTest {
		return nil
	}
	if len(gf.Types) > 1 {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "F2",
			Level:   "ERROR",
			Message: "file contains multiple types; expected 1 file 1 type",
		}}
	}
	return nil
}
