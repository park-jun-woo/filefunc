//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_C1
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/parse"
)

// C1
func TestMutest_C1(t *testing.T) {
	cb, err := parse.ParseCodebook("testdata/codebook_empty_required.yaml")
	if err != nil {
		t.Fatal(err)
	}
	expectViolation(t, CheckCodebookRequiredKeys(cb), "C1")
}
