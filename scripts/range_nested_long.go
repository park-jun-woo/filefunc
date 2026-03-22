// range_nested_long: range-nested에서 순수 줄수 11줄 이상 블록 위치 출력
//
// Usage: go run scripts/range_nested_long.go [dir]
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type entry struct {
	file      string
	fn        string
	line      int
	totalBody int
	innerCtrl int
	pureBody  int
}

var fset *token.FileSet

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	var results []entry
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.Contains(path, "vendor/") || strings.Contains(path, "_test.go") {
			return nil
		}
		collectFile(path, dir, &results)
		return nil
	})

	sort.Slice(results, func(i, j int) bool { return results[i].pureBody > results[j].pureBody })

	fmt.Printf("%-60s %-35s %5s %5s %5s %5s\n", "FILE", "FUNC", "LINE", "TOTAL", "CTRL", "PURE")
	fmt.Println(strings.Repeat("-", 120))
	for _, r := range results {
		if r.pureBody < 11 {
			break
		}
		fmt.Printf("%-60s %-35s %5d %5d %5d %5d\n", r.file, r.fn, r.line, r.totalBody, r.innerCtrl, r.pureBody)
	}
}

func collectFile(path, dir string, results *[]entry) {
	fset = token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return
	}
	rel := path
	if strings.HasPrefix(rel, dir+"/") {
		rel = rel[len(dir)+1:]
	}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		fnName := fn.Name.Name
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			t := fn.Recv.List[0].Type
			if star, ok := t.(*ast.StarExpr); ok {
				t = star.X
			}
			if ident, ok := t.(*ast.Ident); ok {
				fnName = ident.Name + "." + fnName
			}
		}
		for _, stmt := range fn.Body.List {
			rs, ok := stmt.(*ast.RangeStmt)
			if !ok {
				continue
			}
			if !hasInnerControl(rs.Body.List) {
				continue
			}
			total := bodyLines(rs.Body)
			inner := sumInnerControlLines(rs.Body.List)
			pure := total - inner
			*results = append(*results, entry{rel, fnName, fset.Position(rs.Pos()).Line, total, inner, pure})
		}
	}
}

func hasInnerControl(stmts []ast.Stmt) bool {
	for _, stmt := range stmts {
		switch stmt.(type) {
		case *ast.ForStmt, *ast.RangeStmt, *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt:
			return true
		}
	}
	return false
}

func sumInnerControlLines(stmts []ast.Stmt) int {
	total := 0
	for _, stmt := range stmts {
		total += controlLines(stmt)
	}
	return total
}

func controlLines(stmt ast.Stmt) int {
	switch st := stmt.(type) {
	case *ast.ForStmt:
		return stmtSpan(st)
	case *ast.RangeStmt:
		return stmtSpan(st)
	case *ast.IfStmt:
		return ifChainSpan(st)
	case *ast.SwitchStmt:
		return stmtSpan(st)
	case *ast.TypeSwitchStmt:
		return stmtSpan(st)
	}
	return 0
}

func stmtSpan(stmt ast.Stmt) int {
	start := fset.Position(stmt.Pos()).Line
	end := fset.Position(stmt.End()).Line
	return end - start + 1
}

func ifChainSpan(stmt *ast.IfStmt) int {
	start := fset.Position(stmt.Pos()).Line
	end := fset.Position(stmt.End()).Line
	return end - start + 1
}

func bodyLines(block *ast.BlockStmt) int {
	if block == nil {
		return 0
	}
	start := fset.Position(block.Lbrace).Line
	end := fset.Position(block.Rbrace).Line
	lines := end - start - 1
	if lines < 0 {
		return 0
	}
	return lines
}
