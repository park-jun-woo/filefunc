//ff:func feature=validate type=rule control=sequence
//ff:what F3 toulmin rule — 파일당 method 1개 위반 여부를 bool로 반환
package validate

// RuleF3 returns true if the file violates F3 (multiple methods).
func RuleF3(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(gf.Methods) > 1
}
