//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A8
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// A8
func TestMutest_A8(t *testing.T) {
	cb := &model.Codebook{
		Required: map[string]map[string]string{
			"feature": {"validate": "", "parse": ""},
			"type":    {"rule": "", "parser": ""},
		},
	}
	expectViolation(t, ruleViolations(InCodebook, mustParse(t, "testdata/missing_required_key.go"), cb, backingA8), "A8")
}
