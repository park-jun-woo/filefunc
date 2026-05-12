//ff:func feature=parse type=model control=sequence
//ff:what GoFile의 변수 목록을 반환하는 SourceFile 구현 메서드
package model

// GetVars returns the list of top-level var names.
func (gf *GoFile) GetVars() []string { return gf.Vars }
