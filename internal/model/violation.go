//ff:type feature=validate type=model
//ff:what 검증 위반 결과를 담는 구조체
package model

// Violation represents a single rule violation found during validation.
type Violation struct {
	File    string // file path where the violation occurred
	Rule    string // rule ID (F1, Q1, A1, etc.)
	Level   string // ERROR or WARNING
	Message string // human-readable description
}
