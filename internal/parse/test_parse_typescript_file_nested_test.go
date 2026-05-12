//ff:func feature=parse type=util control=sequence
//ff:what test: TestParseTypeScriptFileNested
package parse

import (
	"os/exec"
	"testing"
)

func TestParseTypeScriptFileNested(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not found, skipping TypeScript parse test")
	}

	root := findTsProjectRoot()
	tf, err := ParseTypeScriptFile("testdata/nested_control.ts", root)
	if err != nil {
		t.Fatalf("ParseTypeScriptFile failed: %v", err)
	}

	if tf.MaxDepth != 2 {
		t.Errorf("MaxDepth = %d, want 2", tf.MaxDepth)
	}

	if tf.Control != "iteration" {
		t.Errorf("Control = %q, want %q", tf.Control, "iteration")
	}

	if !tf.HasLoopAtDepth1 {
		t.Error("HasLoopAtDepth1 = false, want true")
	}
}
