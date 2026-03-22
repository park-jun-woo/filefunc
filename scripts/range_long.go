// range_long: range body 21줄 이상 블록의 위치 출력
//
// Usage: go run scripts/range_long.go [dir]
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
	file  string
	fn    string
	line  int
	lines int
}

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
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
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
				collectRange(fset, stmt, path, fnName, &results)
			}
		}
		return nil
	})
	sort.Slice(results, func(i, j int) bool { return results[i].lines > results[j].lines })
	fmt.Printf("%-60s %-40s %6s %5s\n", "FILE", "FUNC", "LINE", "BODY")
	fmt.Println(strings.Repeat("-", 115))
	for _, r := range results {
		if r.lines < 21 {
			break
		}
		rel := r.file
		if strings.HasPrefix(rel, dir+"/") {
			rel = rel[len(dir)+1:]
		}
		fmt.Printf("%-60s %-40s %6d %5d\n", rel, r.fn, r.line, r.lines)
	}
	fmt.Println(strings.Repeat("-", 115))
	fmt.Printf("Total: %d range blocks with body >= 21 lines\n", countOver(results, 21))
}

func collectRange(fset *token.FileSet, stmt ast.Stmt, file, fn string, results *[]entry) {
	s, ok := stmt.(*ast.RangeStmt)
	if !ok {
		return
	}
	lines := bodyLines(fset, s.Body)
	line := fset.Position(s.Pos()).Line
	*results = append(*results, entry{file, fn, line, lines})
}

func bodyLines(fset *token.FileSet, block *ast.BlockStmt) int {
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

func countOver(results []entry, threshold int) int {
	count := 0
	for _, r := range results {
		if r.lines >= threshold {
			count++
		}
	}
	return count
}
