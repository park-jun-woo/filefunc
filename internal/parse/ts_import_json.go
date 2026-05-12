//ff:type feature=parse type=model
//ff:what ts_ast.js JSON의 import 항목 구조체
package parse

// TsImportJSON is the JSON representation of an import from ts_ast.js.
type TsImportJSON struct {
	Module string   `json:"module"`
	Names  []string `json:"names"`
}
