//ff:func feature=validate type=rule control=sequence
//ff:what A15 toulmin rule — control=iteration이면 dimension= 필수 여부
package validate

// RuleA15 returns true if the file violates A15.
func RuleA15(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckDimensionRequired(gf)) > 0
}
