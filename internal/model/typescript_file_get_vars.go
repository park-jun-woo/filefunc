//ff:func feature=parse type=model control=sequence
//ff:what TypeScriptFile의 모듈 레벨 변수 목록을 반환하는 SourceFile 구현 메서드
package model

// GetVars returns the list of module-level UPPER_CASE const names.
func (tf *TypeScriptFile) GetVars() []string { return tf.Vars }
