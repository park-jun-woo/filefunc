//ff:func feature=validate type=rule control=sequence
//ff:what A7 TypeScript 전용 — //ff:checked 해시와 TypeScriptFile.BodyHash 불일치 시 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckedHashMatchTypeScript returns (true, []model.Violation) if the TypeScript file
// has a checked hash annotation that does not match TypeScriptFile.BodyHash.
func CheckedHashMatchTypeScript(claim any, ground any, backing any) (bool, any) {
	g := ground.(*ValidateGround)
	if !g.HasChecked {
		return false, nil
	}
	sf := g.File
	tf, ok := sf.(*model.TypeScriptFile)
	if !ok {
		return false, nil
	}
	ann := tf.GetAnnotation()
	if ann == nil || len(ann.Checked) == 0 {
		return false, nil
	}

	savedHash := ann.Checked["hash"]
	if savedHash == "" {
		return false, nil
	}

	if savedHash != tf.BodyHash {
		return true, []model.Violation{{
			File:    tf.GetPath(),
			Rule:    "A7",
			Level:   "ERROR",
			Message: fmt.Sprintf("checked hash mismatch: saved=%s current=%s (run filefunc llmc to re-verify)", savedHash, tf.BodyHash),
		}}
	}
	return false, nil
}
