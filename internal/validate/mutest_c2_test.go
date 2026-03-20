//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_C2
package validate

import "testing"

// C2
func TestMutest_C2(t *testing.T) {
	expectViolation(t, CheckCodebookDuplicates("testdata/codebook_duplicate_key.yaml"), "C2")
}
