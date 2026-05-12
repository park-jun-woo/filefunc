//ff:func feature=parse type=model control=sequence
//ff:what GoFile의 함수 목록을 반환하는 SourceFile 구현 메서드
package model

// GetFuncs returns the list of top-level func names.
func (gf *GoFile) GetFuncs() []string { return gf.Funcs }
