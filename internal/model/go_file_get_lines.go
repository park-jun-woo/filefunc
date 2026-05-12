//ff:func feature=parse type=model control=sequence
//ff:what GoFile의 총 라인 수를 반환하는 SourceFile 구현 메서드
package model

// GetLines returns the total line count.
func (gf *GoFile) GetLines() int { return gf.Lines }
