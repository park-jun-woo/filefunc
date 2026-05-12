//ff:type feature=parse type=model
//ff:what 언어 중립적 소스 파일 인터페이스 — validate, chain, cli, context 레이어의 공통 추상화
package model

// SourceFile is a language-neutral interface for a parsed source file.
type SourceFile interface {
	GetPath() string
	GetLang() string
	GetPackage() string
	GetFuncs() []string
	GetTypes() []string
	GetMethods() []string
	GetHasInit() bool
	GetVars() []string
	GetAnnotation() *Annotation
	GetLines() int
	GetMaxDepth() int
	IsTestFile() bool
}
