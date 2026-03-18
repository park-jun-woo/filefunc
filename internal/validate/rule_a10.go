//ff:func feature=validate type=rule control=sequence
//ff:what A10 toulmin rule — control=selection인데 switch 없는지 여부
package validate

// RuleA10 returns true if the file violates A10.
func RuleA10(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckControlSelection(gf)) > 0
}
