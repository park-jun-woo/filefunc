//ff:func feature=validate type=rule control=sequence
//ff:what A6 toulmin rule — 어노테이션이 파일 최상단에 위치하는지 여부
package validate

// RuleA6 returns true if the file violates A6 (annotation not at top).
func RuleA6(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	return len(CheckAnnotationPosition(gf)) > 0
}
