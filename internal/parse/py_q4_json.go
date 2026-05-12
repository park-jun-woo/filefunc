//ff:type feature=parse type=model
//ff:what py_ast.py JSON의 Q4 위반 항목 구조체
package parse

// PyQ4JSON is the JSON representation of a Q4 result from py_ast.py.
type PyQ4JSON struct {
	FuncName  string `json:"func_name"`
	StmtType  string `json:"stmt_type"`
	PureLines int    `json:"pure_lines"`
	Line      int    `json:"line"`
}
