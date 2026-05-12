//ff:func feature=validate type=rule control=sequence
//ff:what A7 toulmin rule — //ff:checked 해시 불일치 시 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// RuleA7 returns (true, []model.Violation) if the file violates A7 (checked hash mismatch).
func CheckedHashMatch(claim any, ground any, backing any) (bool, any) {
	g := ground.(*ValidateGround)
	if !g.HasChecked {
		return false, nil
	}
	sf := g.File
	ann := sf.GetAnnotation()
	if ann == nil || len(ann.Checked) == 0 {
		return false, nil
	}

	savedHash := ann.Checked["hash"]
	if savedHash == "" {
		return false, nil
	}

	currentHash, err := parse.CalcBodyHash(sf.GetPath())
	if err != nil {
		return false, nil
	}

	if savedHash != currentHash {
		return true, []model.Violation{{
			File:    sf.GetPath(),
			Rule:    "A7",
			Level:   "ERROR",
			Message: fmt.Sprintf("checked hash mismatch: saved=%s current=%s (run filefunc llmc to re-verify)", savedHash, currentHash),
		}}
	}
	return false, nil
}
