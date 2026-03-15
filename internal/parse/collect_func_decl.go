//ff:func feature=parse type=parser
//ff:what FuncDecl에서 func/method/init 정보를 GoFile에 수집
//ff:checked llm=gpt-oss:20b hash=19e57b4c
package parse

import (
	"go/ast"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// CollectFuncDecl extracts func/method/init info from a FuncDecl into GoFile.
func CollectFuncDecl(gf *model.GoFile, d *ast.FuncDecl) {
	if d.Name.Name == "init" {
		gf.HasInit = true
		return
	}
	if d.Recv != nil {
		recvType := ReceiverTypeName(d.Recv)
		gf.Methods = append(gf.Methods, recvType+"."+d.Name.Name)
		return
	}
	gf.Funcs = append(gf.Funcs, d.Name.Name)
}
