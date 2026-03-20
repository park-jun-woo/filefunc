//ff:func feature=parse type=util control=sequence
//ff:what test: TestMatchPattern_Glob
package walk

import "testing"

func TestMatchPattern_Glob(t *testing.T) {
	if !matchPattern("test.go", "test.go", false, "*.go") {
		t.Error("expected *.go to match test.go")
	}
	if matchPattern("test.txt", "test.txt", false, "*.go") {
		t.Error("expected *.go not to match test.txt")
	}
}
