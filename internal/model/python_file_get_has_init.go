//ff:func feature=parse type=model control=sequence
//ff:what PythonFile의 init 함수 유무를 반환하는 SourceFile 구현 메서드 (Python은 항상 false)
package model

// GetHasInit returns false — Python has no Go-style init().
func (pf *PythonFile) GetHasInit() bool { return false }
