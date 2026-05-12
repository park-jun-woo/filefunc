//ff:func feature=parse type=model control=sequence
//ff:what PythonFile의 클래스 목록을 반환하는 SourceFile 구현 메서드
package model

// GetTypes returns the list of class names.
func (pf *PythonFile) GetTypes() []string { return pf.Classes }
