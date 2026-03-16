//ff:func feature=validate type=rule control=sequence
//ff:what A12: control=sequence인데 depth 1에 switch/loop 존재하면 ERROR
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// CheckControlSequence checks A12: control=sequence must not have
// switch or loop at depth 1.
func CheckControlSequence(gf *model.GoFile) []model.Violation {
	if gf.IsTest || len(gf.Funcs) == 0 || gf.Annotation == nil {
		return nil
	}
	control := gf.Annotation.Func["control"]
	if control != "" && control != "sequence" {
		return nil
	}
	actual := parse.DetectControl(gf.Path)
	if actual != "sequence" {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A12",
			Level:   "ERROR",
			Message: fmt.Sprintf("control=sequence but %s found at depth 1; add control=%s or extract to separate func", actual, actual),
		}}
	}
	return nil
}
