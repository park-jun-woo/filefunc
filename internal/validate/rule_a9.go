//ff:func feature=validate type=rule control=sequence
//ff:what A9 toulmin rule — func 파일의 control= 필수 여부
package validate

// RuleA9 returns true if the file violates A9 (missing control).
func RuleA9(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckControlRequired(gf)) > 0
}
