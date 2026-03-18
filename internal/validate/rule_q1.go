//ff:func feature=validate type=rule control=sequence
//ff:what Q1 toulmin rule — nesting depth 상한 위반 여부를 bool로 반환
package validate

// RuleQ1 returns true if the file violates Q1 (nesting depth exceeds limit).
func RuleQ1(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckNestingDepth(gf)) > 0
}
