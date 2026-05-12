//ff:func feature=parse type=model control=sequence
//ff:what TypeScriptFile의 언어("typescript")를 반환하는 SourceFile 구현 메서드
package model

// GetLang returns "typescript" for TypeScript source files.
func (tf *TypeScriptFile) GetLang() string { return "typescript" }
