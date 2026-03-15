//ff:func feature=parse type=parser
//ff:what GenDecl에서 type/var 정보를 GoFile에 수집
package parse

import (
	"go/ast"
	"go/token"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CollectGenDecl extracts type/var info from a GenDecl into GoFile.
func CollectGenDecl(gf *model.GoFile, d *ast.GenDecl) {
	for _, spec := range d.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if ok {
			gf.Types = append(gf.Types, ts.Name.Name)
			continue
		}
		vs, ok := spec.(*ast.ValueSpec)
		if ok && d.Tok == token.VAR {
			gf.Vars = append(gf.Vars, CollectVarNames(vs)...)
		}
	}
}
