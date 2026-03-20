//ff:func feature=cli type=util control=sequence
//ff:what test: TestBuildPrompt
package llm

import (
	"strings"
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	p := BuildPrompt("validates file structure", "func Check() {}")
	if !strings.Contains(p, "validates file structure") {
		t.Error("missing what")
	}
	if !strings.Contains(p, "func Check() {}") {
		t.Error("missing body")
	}
	if !strings.Contains(p, "0.0") || !strings.Contains(p, "1.0") {
		t.Error("missing score range")
	}
}
