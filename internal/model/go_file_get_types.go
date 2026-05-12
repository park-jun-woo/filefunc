//ff:func feature=parse type=model control=sequence
//ff:what GoFile의 타입 목록을 반환하는 SourceFile 구현 메서드
package model

// GetTypes returns the list of type names.
func (gf *GoFile) GetTypes() []string { return gf.Types }
