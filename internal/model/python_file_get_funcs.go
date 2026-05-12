//ff:func feature=parse type=model control=sequence
//ff:what PythonFile의 함수 목록을 반환하는 SourceFile 구현 메서드
package model

// GetFuncs returns the list of module-level function names.
func (pf *PythonFile) GetFuncs() []string { return pf.Funcs }
