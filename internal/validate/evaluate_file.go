//ff:func feature=validate type=engine control=iteration dimension=1
//ff:what 단일 파일에 대해 toulmin graph를 실행하고 verdict > 0인 룰의 위반을 수집
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// evaluateFile runs the defeats graph on a single file and collects violations.
func evaluateFile(gf *model.GoFile, cb *model.Codebook, ground *ValidateGround) []model.Violation {
	results := ValidateGraph.Evaluate(gf.Path, ground)
	var violations []model.Violation
	for _, r := range results {
		if r.Verdict > 0 {
			violations = append(violations, violationsFor(r.Name, gf, cb)...)
		}
	}
	return violations
}
