//ff:func feature=parse type=model control=sequence
//ff:what GoFile의 init 함수 유무를 반환하는 SourceFile 구현 메서드
package model

// GetHasInit returns whether the file contains func init().
func (gf *GoFile) GetHasInit() bool { return gf.HasInit }
