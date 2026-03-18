//ff:func feature=validate type=rule control=sequence
//ff:what F1 toulmin rule — 파일당 func 1개 위반 여부를 bool로 반환
package validate

// RuleF1 returns true if the file violates F1 (multiple funcs).
// Exception logic (test files, const-only) is handled by defeats graph.
func RuleF1(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(gf.Funcs) > 1
}
