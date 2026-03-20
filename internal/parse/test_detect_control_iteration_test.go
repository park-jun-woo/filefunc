//ff:func feature=parse type=util control=sequence
//ff:what test: TestDetectControl_Iteration
package parse

import (
	"testing"
)

func TestDetectControl_Iteration(t *testing.T) {
	got := DetectControl("testdata/detect_iteration.go")
	if got != "iteration" {
		t.Errorf("DetectControl(iteration) = %q, want %q", got, "iteration")
	}
}
