//ff:func feature=parse type=model control=sequence
//ff:what TypeScriptFile의 파일 경로를 반환하는 SourceFile 구현 메서드
package model

// GetPath returns the file path.
func (tf *TypeScriptFile) GetPath() string { return tf.Path }
