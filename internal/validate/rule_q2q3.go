//ff:func feature=validate type=rule control=sequence
//ff:what Q2/Q3 toulmin rule — func 라인 수 위반 여부를 bool로 반환
package validate

// RuleQ2Q3 returns true if the file violates Q2 or Q3 (func line limits).
func RuleQ2Q3(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckFuncLines(gf)) > 0
}
