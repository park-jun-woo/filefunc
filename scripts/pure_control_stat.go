// pure_control_stat: 제어문 순수 줄수 통계
// - flat: 내부에 제어문이 없는 순수 제어문의 body 줄수
// - nested: 내부에 제어문이 있는 제어문에서 내부 제어문 줄수를 뺀 순수 줄수
// 대상: for, range, if, case (switch 전체가 아닌 case 단위)
//
// Usage: go run scripts/pure_control_stat.go [dir]
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type bucket struct {
	b1_5   int
	b6_10  int
	b11_20 int
	b21_50 int
	b51    int
	total  int
	max    int
}

func (b *bucket) add(lines int) {
	b.total++
	if lines > b.max {
		b.max = lines
	}
	switch {
	case lines <= 5:
		b.b1_5++
	case lines <= 10:
		b.b6_10++
	case lines <= 20:
		b.b11_20++
	case lines <= 50:
		b.b21_50++
	default:
		b.b51++
	}
}

type stats struct {
	forFlat     bucket
	forNested   bucket
	rangeFlat   bucket
	rangeNested bucket
	ifFlat      bucket
	ifNested    bucket
	caseFlat    bucket
	caseNested  bucket
}

var fset *token.FileSet

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	var s stats
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.Contains(path, "vendor/") || strings.Contains(path, "_test.go") {
			return nil
		}
		collectFile(path, &s)
		return nil
	})

	fmt.Printf("%-14s %6s  %6s  %6s  %6s  %6s  %6s  %6s\n",
		"KIND", "TOTAL", "1-5", "6-10", "11-20", "21-50", "51+", "MAX")
	fmt.Println(strings.Repeat("-", 70))
	printRow("for-flat", &s.forFlat)
	printRow("for-nested", &s.forNested)
	fmt.Println()
	printRow("range-flat", &s.rangeFlat)
	printRow("range-nested", &s.rangeNested)
	fmt.Println()
	printRow("if-flat", &s.ifFlat)
	printRow("if-nested", &s.ifNested)
	fmt.Println()
	printRow("case-flat", &s.caseFlat)
	printRow("case-nested", &s.caseNested)
}

func printRow(name string, b *bucket) {
	fmt.Printf("%-14s %6d  %6d  %6d  %6d  %6d  %6d  %6d\n",
		name, b.total, b.b1_5, b.b6_10, b.b11_20, b.b21_50, b.b51, b.max)
}

func collectFile(path string, s *stats) {
	fset = token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return
	}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		for _, stmt := range fn.Body.List {
			collectDepth1(stmt, s)
		}
	}
}

func collectDepth1(stmt ast.Stmt, s *stats) {
	switch st := stmt.(type) {
	case *ast.ForStmt:
		totalLines := bodyLines(st.Body)
		innerLines := sumInnerControlLines(st.Body.List)
		hasInner := hasInnerControl(st.Body.List)
		if hasInner {
			s.forNested.add(totalLines - innerLines)
		} else {
			s.forFlat.add(totalLines)
		}

	case *ast.RangeStmt:
		totalLines := bodyLines(st.Body)
		innerLines := sumInnerControlLines(st.Body.List)
		hasInner := hasInnerControl(st.Body.List)
		if hasInner {
			s.rangeNested.add(totalLines - innerLines)
		} else {
			s.rangeFlat.add(totalLines)
		}

	case *ast.IfStmt:
		totalLines := bodyLines(st.Body)
		innerLines := sumInnerControlLines(st.Body.List)
		hasInner := hasInnerControl(st.Body.List)
		if hasInner {
			s.ifNested.add(totalLines - innerLines)
		} else {
			s.ifFlat.add(totalLines)
		}

	case *ast.SwitchStmt:
		collectCases(st.Body, s)

	case *ast.TypeSwitchStmt:
		collectCases(st.Body, s)
	}
}

func collectCases(body *ast.BlockStmt, s *stats) {
	if body == nil {
		return
	}
	for _, stmt := range body.List {
		cc, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}
		totalLines := caseBodyLines(cc)
		innerLines := sumInnerControlLines(cc.Body)
		hasInner := hasInnerControl(cc.Body)
		if hasInner {
			s.caseNested.add(totalLines - innerLines)
		} else {
			s.caseFlat.add(totalLines)
		}
	}
}

func caseBodyLines(cc *ast.CaseClause) int {
	if len(cc.Body) == 0 {
		return 0
	}
	start := fset.Position(cc.Colon).Line
	last := cc.Body[len(cc.Body)-1]
	end := fset.Position(last.End()).Line
	lines := end - start
	if lines < 0 {
		return 0
	}
	return lines
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
