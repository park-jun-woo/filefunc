//ff:func feature=validate type=rule control=sequence
//ff:what A16 toulmin rule — dimension= 값이 양의 정수인지 여부
package validate

// RuleA16 returns true if the file violates A16.
func RuleA16(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckDimensionValue(gf)) > 0
}
