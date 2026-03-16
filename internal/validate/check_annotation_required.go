//ff:func feature=validate type=rule control=sequence
//ff:what A1: func 파일은 //ff:func, type 파일은 //ff:type 필수 검증
//ff:checked llm=gpt-oss:20b hash=a1630e03
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckAnnotationRequired checks A1: files with funcs must have //ff:func,
// files with types must have //ff:type.
func CheckAnnotationRequired(gf *model.GoFile) []model.Violation {
	if gf.IsTest {
		return nil
	}

	hasFuncs := len(gf.Funcs) > 0
	hasTypes := len(gf.Types) > 0 && !hasFuncs
	ann := gf.Annotation

	if hasFuncs && (ann == nil || len(ann.Func) == 0) {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A1",
			Level:   "ERROR",
			Message: "file with func must have //ff:func annotation",
		}}
	}

	if hasTypes && (ann == nil || len(ann.Type) == 0) {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "A1",
			Level:   "ERROR",
			Message: "file with type must have //ff:type annotation",
		}}
	}

	return nil
}
