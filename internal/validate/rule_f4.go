//ff:func feature=validate type=rule control=sequence
//ff:what F4 toulmin rule — init()만 단독 존재하는 파일 위반 여부를 bool로 반환
package validate

// RuleF4 returns true if the file violates F4 (init-only file).
func RuleF4(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckInitStandalone(gf)) > 0
}
