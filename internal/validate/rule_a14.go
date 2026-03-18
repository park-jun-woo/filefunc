//ff:func feature=validate type=rule control=sequence
//ff:what A14 toulmin rule — control=iteration인데 switch 존재 여부
package validate

// RuleA14 returns true if the file violates A14.
func RuleA14(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckControlIterationNoSwitch(gf)) > 0
}
