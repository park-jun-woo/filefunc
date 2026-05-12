//ff:func feature=parse type=model control=sequence
//ff:what PythonFile의 메서드 목록을 반환하는 SourceFile 구현 메서드
package model

// GetMethods returns the list of method names (ClassName.method).
func (pf *PythonFile) GetMethods() []string { return pf.Methods }
