//ff:func feature=validate type=util control=sequence
//ff:what test: TestMutest_A7
package validate

import "testing"

// A7
func TestMutest_A7(t *testing.T) {
	expectViolation(t, ruleViolations(CheckedHashMatch, mustParse(t, "testdata/checked_hash_mismatch.go"), nil, nil), "A7")
}
