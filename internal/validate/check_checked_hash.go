//ff:func feature=validate type=rule control=sequence
//ff:what A7: //ff:checked 해시와 현재 body 해시를 대조하여 불일치 시 ERROR
//ff:checked llm=gpt-oss:20b hash=d7b8f892
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// CheckCheckedHash checks A7: if //ff:checked exists, its hash must match
// the current func body hash. Mismatch means the body changed after LLM verification.
func CheckCheckedHash(gf *model.GoFile) []model.Violation {
	if gf.Annotation == nil || len(gf.Annotation.Checked) == 0 {
		return nil
	}

	savedHash := gf.Annotation.Checked["hash"]
	if savedHash == "" {
		return nil
	}

	currentHash, err := parse.CalcBodyHash(gf.Path)
	if err != nil {
		return nil
	}

	if savedHash != currentHash {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A7",
			Level:   "ERROR",
			Message: fmt.Sprintf("checked hash mismatch: saved=%s current=%s (run filefunc llmc to re-verify)", savedHash, currentHash),
		}}
	}
	return nil
}
