//ff:func feature=parse type=model control=sequence
//ff:what PythonFile의 언어("python")를 반환하는 SourceFile 구현 메서드
package model

// GetLang returns "python" for Python source files.
func (pf *PythonFile) GetLang() string { return "python" }
