//ff:func feature=validate type=rule control=sequence
//ff:what A7 Python 전용 — # ff:checked 해시와 PythonFile.BodyHash 불일치 시 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckedHashMatchPython returns (true, []model.Violation) if the Python file
// has a checked hash annotation that does not match PythonFile.BodyHash.
func CheckedHashMatchPython(claim any, ground any, backing any) (bool, any) {
	g := ground.(*ValidateGround)
	if !g.HasChecked {
		return false, nil
	}
	sf := g.File
	pf, ok := sf.(*model.PythonFile)
	if !ok {
		return false, nil
	}
	ann := pf.GetAnnotation()
	if ann == nil || len(ann.Checked) == 0 {
		return false, nil
	}

	savedHash := ann.Checked["hash"]
	if savedHash == "" {
		return false, nil
	}

	if savedHash != pf.BodyHash {
		return true, []model.Violation{{
			File:    pf.GetPath(),
			Rule:    "A7",
			Level:   "ERROR",
			Message: fmt.Sprintf("checked hash mismatch: saved=%s current=%s (run filefunc llmc to re-verify)", savedHash, pf.BodyHash),
		}}
	}
	return false, nil
}
