//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A2
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// A2
func TestMutest_A2(t *testing.T) {
	cb := &model.Codebook{
		Required: map[string]map[string]string{
			"feature": {"validate": "", "parse": ""},
			"type":    {"rule": "", "parser": ""},
		},
	}
	expectViolation(t, ruleViolations(InCodebook, mustParse(t, "testdata/bad_codebook_value.go"), cb, backingA2), "A2")
}
