//ff:func feature=validate type=rule
//ff:what Q2/Q3: func 라인 수 검증 (1000 ERROR, 100 WARNING)
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
// Q3: WARNING if func exceeds 100 lines.
func CheckFuncLines(gf *model.GoFile) []model.Violation {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, gf.Path, nil, 0)
	if err != nil {
		return nil
	}

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
		} else if lines > 100 {
			violations = append(violations, model.Violation{
				File:    gf.Path,
				Rule:    "Q3",
				Level:   "WARNING",
				Message: fmt.Sprintf("func %s is %d lines; recommended maximum is 100", fd.Name.Name, lines),
			})
		}
	}
	return violations
}
