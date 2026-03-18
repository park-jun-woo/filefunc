//ff:func feature=validate type=command control=iteration dimension=1
//ff:what 모든 검증 룰을 toulmin defeats graph로 실행하고 위반 목록을 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RunAll executes all validation rules via the toulmin defeats graph.
// Verdict > 0 means the warrant holds (violation). Verdict <= 0 means defeated (exception).
func RunAll(files []*model.GoFile, cb *model.Codebook) []model.Violation {
	var violations []model.Violation
	hasChecked := HasAnyChecked(files)
	for _, gf := range files {
		ground := &ValidateGround{File: gf, Codebook: cb, HasChecked: hasChecked}
		violations = append(violations, evaluateFile(gf, cb, ground)...)
	}
	return violations
}
