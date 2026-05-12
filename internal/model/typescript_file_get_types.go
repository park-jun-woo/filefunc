//ff:func feature=parse type=model control=sequence
//ff:what TypeScriptFile의 타입 목록(Classes + Interfaces + TypeAliases 합산)을 반환하는 SourceFile 구현 메서드
package model

// GetTypes returns the combined list of class, interface, and type alias names.
func (tf *TypeScriptFile) GetTypes() []string {
	result := make([]string, 0, len(tf.Classes)+len(tf.Interfaces)+len(tf.TypeAliases))
	result = append(result, tf.Classes...)
	result = append(result, tf.Interfaces...)
	result = append(result, tf.TypeAliases...)
	return result
}
