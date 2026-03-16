//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what 파일 내 모든 함수의 최대 nesting depth 계산
//ff:checked llm=gpt-oss:20b hash=d75547cd
package parse

import "go/ast"

// CalcMaxDepth calculates the maximum nesting depth across all functions in a file.
func CalcMaxDepth(f *ast.File) int {
	maxDepth := 0
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		d := StmtDepth(fd.Body.List, 0)
		if d > maxDepth {
			maxDepth = d
		}
	}
	return maxDepth
}
