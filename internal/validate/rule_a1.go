//ff:func feature=validate type=rule control=sequence
//ff:what A1 toulmin rule — func/type 파일의 //ff:func 또는 //ff:type 필수 위반 시 violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RuleA1 returns (true, []model.Violation) if the file violates A1 (missing annotation).
func RuleA1(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	hasFuncs := len(gf.Funcs) > 0
	hasTypes := len(gf.Types) > 0
	ann := gf.Annotation

	if hasFuncs && (ann == nil || len(ann.Func) == 0) {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A1",
			Level:   "ERROR",
			Message: "file with func must have //ff:func annotation",
		}}
	}

	if hasTypes && (ann == nil || len(ann.Type) == 0) {
		return true, []model.Violation{{
			File:    gf.Path,
			Rule:    "A1",
			Level:   "ERROR",
			Message: "file with type must have //ff:type annotation",
		}}
	}

	return false, nil
}
