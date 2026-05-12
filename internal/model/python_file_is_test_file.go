//ff:func feature=parse type=model control=sequence
//ff:what PythonFile의 테스트 파일 여부를 반환하는 SourceFile 구현 메서드
package model

// IsTestFile returns whether the file is a test file (test_ prefix or _test suffix).
func (pf *PythonFile) IsTestFile() bool { return pf.IsTest }
