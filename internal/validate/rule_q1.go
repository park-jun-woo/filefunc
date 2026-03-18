//ff:func feature=validate type=rule control=sequence
//ff:what Q1 toulmin rule — nesting depth 상한 위반 시 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// RuleQ1 returns (true, []model.Violation) if the file violates Q1 (nesting depth exceeds limit).
func RuleQ1(claim any, ground any) (bool, any) {
	gf := ground.(*ValidateGround).File
	limit := depthLimit(gf)
	if gf.MaxDepth > limit {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "Q1",
			Level:   "ERROR",
			Message: fmt.Sprintf("nesting depth %d exceeds maximum of %d", gf.MaxDepth, limit),
		}}
	}
	return false, nil
}
