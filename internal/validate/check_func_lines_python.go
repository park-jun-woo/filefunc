//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what Q2/Q3 Python 전용 — PythonFile.FuncLines에서 func 라인 수 위반 시 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckFuncLinesPython returns (true, []model.Violation) if the Python file violates Q2 or Q3.
// Q2: func > 1000 lines. Q3: control=sequence func > 100 lines.
func CheckFuncLinesPython(claim any, ground any, backing any) (bool, any) {
	sf := ground.(*ValidateGround).File
	pf, ok := sf.(*model.PythonFile)
	if !ok {
		return false, nil
	}

	q3Limit, q3Applies := Q3Limit(sf)
	var violations []model.Violation

	for name, lines := range pf.FuncLines {
		if lines > 1000 {
			violations = append(violations, model.Violation{
				File:    pf.Path,
				Rule:    "Q2",
				Level:   "ERROR",
				Message: fmt.Sprintf("func %s is %d lines; maximum is 1000", name, lines),
			})
			continue
		}
		if q3Applies && lines > q3Limit {
			violations = append(violations, model.Violation{
				File:    pf.Path,
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
