//ff:func feature=validate type=formatter control=iteration dimension=1
//ff:what 순환 경로를 I1 ERROR violation과 fix advice 문자열로 포맷
package validate

import (
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// formatCycleAdvice formats a cycle path into an I1 violation with fix advice.
// The fix advice recommends moving the last edge's import inside a function body.
func formatCycleAdvice(cycle []string) model.Violation {
	names := make([]string, len(cycle))
	for i, p := range cycle {
		names[i] = filepath.Base(p)
	}

	cyclePath := strings.Join(names, " -> ") + " -> " + names[0]

	lastFile := filepath.Base(cycle[len(cycle)-1])
	firstFile := filepath.Base(cycle[0])
	firstModule := strings.TrimSuffix(firstFile, ".py")

	return model.Violation{
		File:  cycle[0],
		Rule:  "I1",
		Level: "ERROR",
		Message: "circular import -- " + cyclePath +
			"\n  fix: move \"from ." + firstModule + " import ...\" inside function body in " + lastFile,
	}
}
