//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what Q2/Q3 TypeScript 전용 — TypeScriptFile.FuncLines에서 func 라인 수 위반 시 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckFuncLinesTypeScript returns (true, []model.Violation) if the TypeScript file violates Q2 or Q3.
// Q2: func > 1000 lines. Q3: control=sequence func > 100 lines.
func CheckFuncLinesTypeScript(claim any, ground any, backing any) (bool, any) {
	sf := ground.(*ValidateGround).File
	tf, ok := sf.(*model.TypeScriptFile)
	if !ok {
		return false, nil
	}

	q3Limit, q3Applies := Q3Limit(sf)
	var violations []model.Violation

	for name, lines := range tf.FuncLines {
		if lines > 1000 {
			violations = append(violations, model.Violation{
				File:    tf.Path,
				Rule:    "Q2",
				Level:   "ERROR",
				Message: fmt.Sprintf("func %s is %d lines; maximum is 1000", name, lines),
			})
			continue
		}
		if q3Applies && lines > q3Limit {
			violations = append(violations, model.Violation{
				File:    tf.Path,
				Rule:    "Q3",
				Level:   "ERROR",
				Message: fmt.Sprintf("func %s is %d lines; maximum for sequence is %d", name, lines, q3Limit),
			})
		}
	}

	if len(violations) > 0 {
		return true, violations
	}
	return false, nil
}
