//ff:func feature=validate type=rule control=sequence
//ff:what Q1: control과 dimension 기반으로 nesting depth 상한 검증
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckNestingDepth checks Q1: nesting depth must not exceed the limit
// based on control type and dimension.
// sequence=2, selection=2, iteration=dimension+1.
func CheckNestingDepth(gf *model.GoFile) []model.Violation {
	limit := depthLimit(gf)
	if gf.MaxDepth > limit {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "Q1",
			Level:   "ERROR",
			Message: fmt.Sprintf("nesting depth %d exceeds maximum of %d", gf.MaxDepth, limit),
		}}
	}
	return nil
}
