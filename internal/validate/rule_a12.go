//ff:func feature=validate type=rule control=sequence
//ff:what A12 toulmin rule — control=sequence인데 switch/loop 존재 시 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// RuleA12 returns (true, []model.Violation) if the file violates A12.
func RuleA12(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	if len(gf.Funcs) == 0 || gf.Annotation == nil {
		return false, nil
	}
	control := gf.Annotation.Func["control"]
	if control != "sequence" {
		return false, nil
	}
	actual := parse.DetectControl(gf.Path)
	if actual != "sequence" {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A12",
			Level:   "ERROR",
			Message: fmt.Sprintf("control=sequence but %s found at depth 1; add control=%s or extract to separate func", actual, actual),
		}}
	}
	return false, nil
}
