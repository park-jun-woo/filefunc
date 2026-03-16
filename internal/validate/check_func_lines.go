//ff:func feature=validate type=rule control=iteration
//ff:what Q2/Q3: func 라인 수 검증. Q3은 control별 기준 (sequence/iteration 100, selection 300)
package validate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CheckFuncLines checks Q2/Q3: func line count.
// Q2: ERROR if func exceeds 1000 lines.
// Q3: WARNING based on control type (sequence/iteration: 100, selection: 300).
func CheckFuncLines(gf *model.GoFile) []model.Violation {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gf.Path, nil, 0)
	if err != nil {
		return nil
	}

	q3Limit := Q3Limit(gf)
	var violations []model.Violation

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}

		startLine := fset.Position(fd.Pos()).Line
		endLine := fset.Position(fd.End()).Line
		lines := endLine - startLine + 1

		if lines > 1000 {
			violations = append(violations, model.Violation{
				File:    gf.Path,
				Rule:    "Q2",
				Level:   "ERROR",
				Message: fmt.Sprintf("func %s is %d lines; maximum is 1000", fd.Name.Name, lines),
			})
			continue
		}

		if lines <= q3Limit {
			continue
		}

		msg := fmt.Sprintf("func %s is %d lines; recommended maximum is %d", fd.Name.Name, lines, q3Limit)
		if HasBacktick(fd) {
			msg += "\n  hint: backtick string detected — consider extracting to a separate template file"
		}
		violations = append(violations, model.Violation{
			File:    gf.Path,
			Rule:    "Q3",
			Level:   "WARNING",
			Message: msg,
		})
	}
	return violations
}
