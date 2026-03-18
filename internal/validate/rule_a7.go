//ff:func feature=validate type=rule control=sequence
//ff:what A7 toulmin rule — //ff:checked 해시 불일치 여부
package validate

// RuleA7 returns true if the file violates A7 (checked hash mismatch).
func RuleA7(claim any, ground any) bool {
	g := ground.(*ValidateGround)
	if !g.HasChecked {
		return false
	}
	return len(CheckCheckedHash(g.File)) > 0
}
