//ff:func feature=parse type=model control=sequence
//ff:what PythonFile의 모듈명을 반환하는 SourceFile 구현 메서드
package model

// GetPackage returns the module name.
func (pf *PythonFile) GetPackage() string { return pf.Module }
