//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what func signature+body의 SHA-256 해시 앞 8자리 계산
//ff:checked llm=gpt-oss:20b hash=f72ddf88
package parse

import (
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// CalcBodyHash calculates the SHA-256 hash (first 8 hex chars) of the main
// func's signature and body in a Go source file.
func CalcBodyHash(path string) (string, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, 0)
	if err != nil {
		return "", err
	}

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil || fd.Name.Name == "init" {
			continue
		}
		start := fset.Position(fd.Pos()).Offset
		end := fset.Position(fd.End()).Offset
		h := sha256.Sum256(src[start:end])
		return fmt.Sprintf("%x", h[:4]), nil
	}

	return "", fmt.Errorf("no func found in %s", path)
}
