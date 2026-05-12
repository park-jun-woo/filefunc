//ff:func feature=parse type=model control=sequence
//ff:what PythonFile의 파일 경로를 반환하는 SourceFile 구현 메서드
package model

// GetPath returns the file path.
func (pf *PythonFile) GetPath() string { return pf.Path }
