//ff:func feature=parse type=parser
//ff:what 메서드 리시버의 타입명을 추출
package parse

import "go/ast"

// ReceiverTypeName extracts the type name from a method receiver field list.
func ReceiverTypeName(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) == 0 {
		return ""
	}
	t := recv.List[0].Type
	if star, ok := t.(*ast.StarExpr); ok {
		t = star.X
	}
	if ident, ok := t.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}
