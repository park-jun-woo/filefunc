//ff:func feature=validate type=command
//ff:what 모든 검증 룰을 실행하고 위반 목록을 반환
//ff:calls CheckAnnotationPosition, CheckAnnotationRequired, CheckCheckedHash, CheckCodebookValues, CheckFuncLines, CheckInitStandalone, CheckNestingDepth, CheckOneFileOneFunc, CheckOneFileOneMethod, CheckOneFileOneType, CheckWhatRequired, HasAnyChecked
//ff:uses Codebook, GoFile, Violation
//ff:checked llm=gpt-oss:20b hash=b5bdcebc
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RunAll executes all validation rules against parsed Go files and returns violations.
func RunAll(files []*model.GoFile, cb *model.Codebook) []model.Violation {
	var violations []model.Violation
	hasChecked := HasAnyChecked(files)
	for _, gf := range files {
		violations = append(violations, CheckOneFileOneFunc(gf)...)
		violations = append(violations, CheckOneFileOneType(gf)...)
		violations = append(violations, CheckOneFileOneMethod(gf)...)
		violations = append(violations, CheckInitStandalone(gf)...)
		violations = append(violations, CheckNestingDepth(gf)...)
		violations = append(violations, CheckFuncLines(gf)...)
		violations = append(violations, CheckAnnotationRequired(gf)...)
		violations = append(violations, CheckCodebookValues(gf, cb)...)
		violations = append(violations, CheckRequiredKeysInAnnotation(gf, cb)...)
		violations = append(violations, CheckWhatRequired(gf)...)
		violations = append(violations, CheckAnnotationPosition(gf)...)
		if hasChecked {
			violations = append(violations, CheckCheckedHash(gf)...)
		}
	}
	return violations
}
