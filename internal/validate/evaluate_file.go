//ff:func feature=validate type=engine control=iteration dimension=1
//ff:what 단일 파일에 대해 toulmin graph를 실행하고 verdict > 0인 룰의 evidence에서 위반을 수집
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// evaluateFile runs the defeats graph on a single file and collects violations
// from the evidence returned by each rule function.
func evaluateFile(gf *model.GoFile, cb *model.Codebook, ground *ValidateGround) []model.Violation {
	results, err := ValidateGraph.Evaluate(gf.Path, ground)
	if err != nil {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "EVAL",
			Level:   "ERROR",
			Message: err.Error(),
		}}
	}
	var violations []model.Violation
	for _, r := range results {
		vs, ok := r.Evidence.([]model.Violation)
		if r.Verdict > 0 && ok {
			violations = append(violations, vs...)
		}
	}
	return violations
}
