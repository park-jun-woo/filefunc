//ff:type feature=parse type=model
//ff:what py_ast.py JSON의 import 항목 구조체
package parse

// PyImportJSON is the JSON representation of an import from py_ast.py.
type PyImportJSON struct {
	Module string   `json:"module"`
	Names  []string `json:"names"`
}
