//ff:func feature=validate type=rule
//ff:what Q1: nesting depth 2 초과 검증
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckNestingDepth checks Q1: nesting depth must not exceed 2.
func CheckNestingDepth(gf *model.GoFile) []model.Violation {
	if gf.MaxDepth > 2 {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "Q1",
			Level:   "ERROR",
			Message: fmt.Sprintf("nesting depth %d exceeds maximum of 2", gf.MaxDepth),
		}}
	}
	return nil
}
