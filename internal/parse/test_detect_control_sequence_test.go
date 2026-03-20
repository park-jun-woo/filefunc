//ff:func feature=parse type=util control=sequence
//ff:what test: TestDetectControl_Sequence
package parse

import (
	"testing"
)

func TestDetectControl_Sequence(t *testing.T) {
	got := DetectControl("testdata/detect_sequence.go")
	if got != "sequence" {
		t.Errorf("DetectControl(sequence) = %q, want %q", got, "sequence")
	}
}
