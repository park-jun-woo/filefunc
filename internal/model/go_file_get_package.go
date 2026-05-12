//ff:func feature=parse type=model control=sequence
//ff:what GoFile의 패키지명을 반환하는 SourceFile 구현 메서드
package model

// GetPackage returns the package name.
func (gf *GoFile) GetPackage() string { return gf.Package }
