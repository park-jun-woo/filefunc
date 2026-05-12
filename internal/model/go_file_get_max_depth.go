//ff:func feature=parse type=model control=sequence
//ff:what GoFile의 최대 네스팅 깊이를 반환하는 SourceFile 구현 메서드
package model

// GetMaxDepth returns the maximum nesting depth.
func (gf *GoFile) GetMaxDepth() int { return gf.MaxDepth }
