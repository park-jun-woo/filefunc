//ff:func feature=parse type=model control=sequence
//ff:what GoFile의 테스트 파일 여부를 반환하는 SourceFile 구현 메서드
package model

// IsTestFile returns whether the file is a _test.go file.
func (gf *GoFile) IsTestFile() bool { return gf.IsTest }
