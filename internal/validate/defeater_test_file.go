//ff:func feature=validate type=rule control=sequence
//ff:what F5 defeater — test 파일은 F1/F2/F3/A 룰에서 예외
//ff:why test files conventionally group multiple test funcs and helpers
package validate

// DefeaterTestFile returns true if the file is a test file.
// Used as a defeater in the defeats graph against F1, F2, F3, and annotation rules.
func DefeaterTestFile(claim any, ground any, backing any) (bool, any) {
	return ground.(*ValidateGround).File.IsTest, nil
}
