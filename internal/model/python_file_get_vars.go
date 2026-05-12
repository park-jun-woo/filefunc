//ff:func feature=parse type=model control=sequence
//ff:what PythonFile의 모듈 레벨 변수 목록을 반환하는 SourceFile 구현 메서드
package model

// GetVars returns the list of module-level uppercase variable names.
func (pf *PythonFile) GetVars() []string { return pf.Vars }
