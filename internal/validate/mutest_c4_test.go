//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_C4
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

// C4
func TestMutest_C4(t *testing.T) {
	cb, err := parse.ParseCodebook("testdata/codebook_no_description.yaml")
	if err != nil {
		t.Fatal(err)
	}
	expectViolation(t, CheckCodebookDescription(cb), "C4")
}
