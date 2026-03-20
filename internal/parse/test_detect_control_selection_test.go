//ff:func feature=parse type=util control=sequence
//ff:what test: TestDetectControl_Selection
package parse

import (
	"testing"
)

func TestDetectControl_Selection(t *testing.T) {
	got := DetectControl("testdata/detect_selection.go")
	if got != "selection" {
		t.Errorf("DetectControl(selection) = %q, want %q", got, "selection")
	}
}
