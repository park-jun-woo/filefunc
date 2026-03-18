//ff:func feature=validate type=engine control=iteration dimension=1
//ff:what 단일 파일에 대해 toulmin graph를 실행하고 verdict > 0인 룰의 evidence에서 위반을 수집
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// evaluateFile runs the defeats graph on a single file and collects violations
// from the evidence returned by each rule function.
func evaluateFile(gf *model.GoFile, cb *model.Codebook, ground *ValidateGround) []model.Violation {
	results := ValidateGraph.Evaluate(gf.Path, ground)
	var violations []model.Violation
	for _, r := range results {
		if r.Verdict > 0 && r.Evidence != nil {
			violations = append(violations, r.Evidence.([]model.Violation)...)
		}
	}
	return violations
}
