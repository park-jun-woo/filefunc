//ff:func feature=validate type=util control=sequence
//ff:what test: expectNoViolation
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func expectNoViolation(t *testing.T, violations []model.Violation) {
	t.Helper()
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}
