//ff:func feature=validate type=rule control=sequence
//ff:what A8 toulmin rule — codebook required 키가 어노테이션에 모두 존재하는지 여부
package validate

// RuleA8 returns true if the file violates A8 (missing required keys).
func RuleA8(claim any, ground any) bool {
	g := ground.(*ValidateGround)
	return len(CheckRequiredKeysInAnnotation(g.File, g.Codebook)) > 0
}
