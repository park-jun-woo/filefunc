//ff:func feature=validate type=rule control=sequence
//ff:what A2 toulmin rule — 어노테이션 값이 코드북에 존재하는지 여부
package validate

// RuleA2 returns true if the file violates A2 (annotation value not in codebook).
func RuleA2(claim any, ground any) bool {
	g := ground.(*ValidateGround)
	return len(CheckCodebookValues(g.File, g.Codebook)) > 0
}
