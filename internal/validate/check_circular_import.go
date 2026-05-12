//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what PythonFile 슬라이스에서 순환 import를 탐지하여 I1 violation 목록 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckCircularImport detects circular imports among Python files.
// It builds an import graph, finds all cycles, and returns I1 violations.
func CheckCircularImport(files []*model.PythonFile, root string) []model.Violation {
	graph := buildImportGraph(files, root)
	cycles := detectCycle(graph)

	var violations []model.Violation
	for _, cycle := range cycles {
		violations = append(violations, formatCycleAdvice(cycle))
	}
	return violations
}
