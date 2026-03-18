//ff:func feature=validate type=rule control=sequence
//ff:what A1 toulmin rule — func/type 파일의 //ff:func 또는 //ff:type 필수 여부
package validate

// RuleA1 returns true if the file violates A1 (missing annotation).
func RuleA1(claim any, ground any) bool {
	gf := ground.(*ValidateGround).File
	ann := gf.Annotation
	hasFuncs := len(gf.Funcs) > 0
	hasTypes := len(gf.Types) > 0 && !hasFuncs
	if hasFuncs && (ann == nil || len(ann.Func) == 0) {
		return true
	}
	if hasTypes && (ann == nil || len(ann.Type) == 0) {
		return true
	}
	return false
}
