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
	gf := g.File
	if gf.Annotation == nil || len(gf.Annotation.Checked) == 0 {
		return false, nil
	}

	savedHash := gf.Annotation.Checked["hash"]
	if savedHash == "" {
		return false, nil
	}

	currentHash, err := parse.CalcBodyHash(gf.Path)
	if err != nil {
		return false, nil
	}

	if savedHash != currentHash {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A7",
			Level:   "ERROR",
			Message: fmt.Sprintf("checked hash mismatch: saved=%s current=%s (run filefunc llmc to re-verify)", savedHash, currentHash),
		}}
	}
	return false, nil
}
