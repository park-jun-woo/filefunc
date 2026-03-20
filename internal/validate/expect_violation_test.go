//ff:func feature=validate type=util control=iteration dimension=1
//ff:what test: expectViolation
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func expectViolation(t *testing.T, violations []model.Violation, rule string) {
	t.Helper()
	if len(violations) == 0 {
		t.Errorf("expected violation %s, got none", rule)
		return
	}
	for _, v := range violations {
		if v.Rule == rule {
			return
		}
	}
	t.Errorf("expected rule %s, got %v", rule, violations)
}
