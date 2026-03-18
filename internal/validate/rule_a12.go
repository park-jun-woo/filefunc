//ff:func feature=validate type=rule control=sequence
//ff:what A12 toulmin rule — control=sequence인데 switch/loop 존재 여부
package validate

// RuleA12 returns true if the file violates A12.
func RuleA12(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckControlSequence(gf)) > 0
}
