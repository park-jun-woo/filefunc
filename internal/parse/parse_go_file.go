//ff:func feature=parse type=parser control=iteration
//ff:what Go 소스 파일을 파싱하여 GoFile 구조체로 변환
//ff:checked llm=gpt-oss:20b hash=b53102b8
package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// ParseGoFile parses a Go source file and extracts its structure.
func ParseGoFile(path string) (*model.GoFile, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	gf := &model.GoFile{
		Path:    path,
		Package: f.Name.Name,
		IsTest:  strings.HasSuffix(path, "_test.go"),
	}

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			CollectFuncDecl(gf, d)
		case *ast.GenDecl:
			CollectGenDecl(gf, d)
		}
	}

	tokFile := fset.File(f.Pos())
	if tokFile != nil {
		gf.Lines = tokFile.LineCount()
	}

	gf.MaxDepth = CalcMaxDepth(f)

	ann, _ := ParseAnnotation(path)
	gf.Annotation = ann

	return gf, nil
}
