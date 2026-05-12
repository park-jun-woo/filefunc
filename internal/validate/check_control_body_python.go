//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what Q4 Python 전용 — PythonFile.Q4Violations에서 PURE 10줄 초과 시 violation 반환
package validate

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckControlBodyPython returns (true, []model.Violation) if a Python control body's
// PURE lines exceed 10. Uses pre-computed Q4Violations from py_ast.py.
func CheckControlBodyPython(claim any, ground any, backing any) (bool, any) {
	sf := ground.(*ValidateGround).File
	pf, ok := sf.(*model.PythonFile)
	if !ok {
		return false, nil
	}

	var violations []model.Violation
	for _, q := range pf.Q4Violations {
		if q.PureLines > 10 {
			violations = append(violations, model.Violation{
				File:  pf.Path,
				Rule:  "Q4",
				Level: "ERROR",
				Message: fmt.Sprintf("func %s: %s body at line %d has %d PURE lines; extract to sequence func (max 10)",
					q.FuncName, q.StmtType, q.Line, q.PureLines),
			})
		}
	}

	if len(violations) > 0 {
		return true, violations
	}
	return false, nil
}
