//ff:func feature=validate type=rule control=sequence
//ff:what A3 toulmin rule — func/type 파일의 //ff:what 필수 여부
package validate

// RuleA3 returns true if the file violates A3 (missing what).
func RuleA3(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	hasFuncs := len(gf.Funcs) > 0 || len(gf.Methods) > 0
	hasTypes := len(gf.Types) > 0
	if !hasFuncs && !hasTypes {
		return false
	}
	return gf.Annotation == nil || gf.Annotation.What == ""
}
