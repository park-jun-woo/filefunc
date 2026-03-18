//ff:func feature=validate type=rule control=sequence
//ff:what F2 toulmin rule — 파일당 type 1개 위반 여부를 bool로 반환
package validate

// RuleF2 returns true if the file violates F2 (multiple exported types).
func RuleF2(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(gf.Types) > 1
}
