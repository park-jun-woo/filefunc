//ff:func feature=validate type=command control=iteration dimension=1
//ff:what 모든 검증 룰을 toulmin defeats graph로 실행하고 위반 목록을 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// RunAll executes all validation rules via the toulmin defeats graph.
// Verdict > 0 means the warrant holds (violation). Verdict <= 0 means defeated (exception).
func RunAll(files []model.SourceFile, cb *model.Codebook) []model.Violation {
	var violations []model.Violation
	hasChecked := HasAnyChecked(files)
	for _, sf := range files {
		ground := &ValidateGround{File: sf, Codebook: cb, HasChecked: hasChecked}
		violations = append(violations, evaluateFile(sf, cb, ground)...)
	}
	return violations
}
