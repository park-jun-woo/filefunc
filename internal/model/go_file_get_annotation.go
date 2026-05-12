//ff:func feature=parse type=model control=sequence
//ff:what GoFile의 어노테이션을 반환하는 SourceFile 구현 메서드
package model

// GetAnnotation returns the parsed //ff: annotation.
func (gf *GoFile) GetAnnotation() *Annotation { return gf.Annotation }
