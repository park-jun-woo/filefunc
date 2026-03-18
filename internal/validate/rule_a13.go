//ff:func feature=validate type=rule control=sequence
//ff:what A13 toulmin rule — control=selection인데 loop 존재 여부
package validate

// RuleA13 returns true if the file violates A13.
func RuleA13(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckControlSelectionNoLoop(gf)) > 0
}
