//ff:func feature=parse type=parser
//ff:what Go AST로 파일의 첫 번째 func(init 제외)의 signature+body를 소스 텍스트로 추출하여 반환
//ff:checked llm=gpt-oss:20b hash=eb39edc3
package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// ExtractFuncSource extracts the source text of the main func (excluding init)
// from a Go source file. Returns empty string if no func found.
func ExtractFuncSource(path string, src []byte) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, 0)
	if err != nil {
		return ""
	}

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Name.Name == "init" {
			continue
		}
		start := fset.Position(fd.Pos()).Offset
		end := fset.Position(fd.End()).Offset
		return string(src[start:end])
	}
	return ""
}
