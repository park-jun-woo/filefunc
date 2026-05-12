//ff:func feature=parse type=model control=sequence
//ff:what TypeScriptFile의 init 함수 유무를 반환하는 SourceFile 구현 메서드 (TypeScript는 항상 false)
package model

// GetHasInit returns false — TypeScript has no Go-style init().
func (tf *TypeScriptFile) GetHasInit() bool { return false }
