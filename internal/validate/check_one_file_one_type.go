//ff:func feature=validate type=rule
//ff:what F2: 파일당 exported type 1개 검증
package validate

import (
	"unicode"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckOneFileOneType checks F2: each file must contain at most one type.
// Exception F6: unexported types used as func-specific parameters are allowed
// alongside the main type or func.
func CheckOneFileOneType(gf *model.GoFile) []model.Violation {
	if gf.IsTest {
		return nil
	}

	exportedTypes := 0
	for _, name := range gf.Types {
		if len(name) > 0 && unicode.IsUpper(rune(name[0])) {
			exportedTypes++
		}
	}

	if exportedTypes > 1 {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "F2",
			Level:   "ERROR",
			Message: "file contains multiple exported types; expected 1 file 1 type",
		}}
	}
	return nil
}
