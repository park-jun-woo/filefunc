//ff:func feature=validate type=rule control=sequence
//ff:what Q1 toulmin rule — nesting depth 상한 위반 시 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// RuleQ1 returns (true, []model.Violation) if the file violates Q1 (nesting depth exceeds limit).
func CheckDepthLimit(claim any, ground any, backing any) (bool, any) {
	sf := ground.(*ValidateGround).File
	limit := depthLimit(sf)
	if sf.GetMaxDepth() > limit {
		return true, []model.Violation{{
			File:    sf.GetPath(),
			Rule:    "Q1",
			Level:   "ERROR",
			Message: fmt.Sprintf("nesting depth %d exceeds maximum of %d", sf.GetMaxDepth(), limit),
		}}
	}
	return false, nil
}
