//ff:func feature=parse type=model control=sequence
//ff:what TypeScriptFile의 테스트 파일 여부를 반환하는 SourceFile 구현 메서드
package model

// IsTestFile returns whether the file is a test file (*.test.ts, *.spec.ts, *.test.tsx, *.spec.tsx).
func (tf *TypeScriptFile) IsTestFile() bool { return tf.IsTest }
