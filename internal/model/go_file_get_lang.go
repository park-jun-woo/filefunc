//ff:func feature=parse type=model control=sequence
//ff:what GoFile의 언어("go")를 반환하는 SourceFile 구현 메서드
package model

// GetLang returns "go" for Go source files.
func (gf *GoFile) GetLang() string { return "go" }
