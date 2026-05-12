//ff:func feature=parse type=util control=sequence
//ff:what test: TestParseTypeScriptFileEmpty
package parse

import (
	"os/exec"
	"testing"
)

func TestParseTypeScriptFileEmpty(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not found, skipping TypeScript parse test")
	}

	root := findTsProjectRoot()
	tf, err := ParseTypeScriptFile("testdata/empty.ts", root)
	if err != nil {
		t.Fatalf("ParseTypeScriptFile failed: %v", err)
	}

	if len(tf.Funcs) != 0 {
		t.Errorf("Funcs = %v, want []", tf.Funcs)
	}

	if len(tf.Classes) != 0 {
		t.Errorf("Classes = %v, want []", tf.Classes)
	}

	if tf.MaxDepth != 0 {
		t.Errorf("MaxDepth = %d, want 0", tf.MaxDepth)
	}

	if tf.Lines != 0 {
		t.Errorf("Lines = %d, want 0", tf.Lines)
	}

	if tf.BodyHash != "" {
		t.Errorf("BodyHash = %q, want empty", tf.BodyHash)
	}
}
