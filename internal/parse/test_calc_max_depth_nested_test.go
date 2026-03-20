//ff:func feature=parse type=util control=sequence
//ff:what test: TestCalcMaxDepth_Nested
package parse

import (
	"testing"
)

func TestCalcMaxDepth_Nested(t *testing.T) {
	gf, err := ParseGoFile("testdata/depth_nested.go")
	if err != nil {
		t.Fatalf("ParseGoFile failed: %v", err)
	}
	// for { if { if {} } } = depth 3
	if gf.MaxDepth != 3 {
		t.Errorf("MaxDepth = %d, want 3", gf.MaxDepth)
	}
}
