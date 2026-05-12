//ff:type feature=validate type=model
//ff:what 제어문 body PURE 줄 수 위반 결과를 담는 구조체
package model

// Q4Result holds a Q4 rule violation: a control body with too many pure lines.
type Q4Result struct {
	FuncName  string
	StmtType  string
	PureLines int
	Line      int
}
