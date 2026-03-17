//ff:func feature=chain type=util control=sequence
//ff:what 패키지명과 함수명을 결합하여 qualified name(pkg.FuncName) 생성
package chain

func qualifiedName(pkg, name string) string {
	return pkg + "." + name
}
