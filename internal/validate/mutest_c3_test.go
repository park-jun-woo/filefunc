//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_C3
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

// C3
func TestMutest_C3(t *testing.T) {
	cb, err := parse.ParseCodebook("testdata/codebook_bad_format.yaml")
	if err != nil {
		t.Fatal(err)
	}
	expectViolation(t, CheckCodebookValueFormat(cb), "C3")
}
