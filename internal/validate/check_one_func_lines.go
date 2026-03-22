//ff:func feature=validate type=util control=sequence
//ff:what 단일 FuncDecl의 Q2/Q3 줄수 위반을 검사하여 violations에 추가
package validate

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func checkOneFuncLines(fset *token.FileSet, path string, fd *ast.FuncDecl, q3Limit int, q3Applies bool, violations []model.Violation) []model.Violation {
	startLine := fset.Position(fd.Pos()).Line
	endLine := fset.Position(fd.End()).Line
	lines := endLine - startLine + 1

	if lines > 1000 {
		violations = append(violations, model.Violation{
			File:    path,
			Rule:    "Q2",
			Level:   "ERROR",
			Message: fmt.Sprintf("func %s is %d lines; maximum is 1000", fd.Name.Name, lines),
		})
		return violations
	}

	if !q3Applies || lines <= q3Limit {
		return violations
	}

	msg := fmt.Sprintf("func %s is %d lines; maximum for sequence is %d", fd.Name.Name, lines, q3Limit)
	if HasBacktick(fd) {
		msg += "\n  hint: backtick string detected — extract to a var-only file (e.g. template_xxx.go) to exempt from Q3"
	}
	violations = append(violations, model.Violation{
		File:    path,
		Rule:    "Q3",
		Level:   "ERROR",
		Message: msg,
	})
	return violations
}
