//ff:func feature=validate type=rule control=sequence
//ff:what A11 toulmin rule — control=iteration인데 loop 없는지 여부
package validate

// RuleA11 returns true if the file violates A11.
func RuleA11(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckControlIteration(gf)) > 0
}
