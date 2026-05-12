//ff:func feature=parse type=util control=sequence
//ff:what test: TestParseTypeScriptFileSimple
package parse

import (
	"os/exec"
	"testing"
)

func TestParseTypeScriptFileSimple(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not found, skipping TypeScript parse test")
	}

	root := findTsProjectRoot()
	tf, err := ParseTypeScriptFile("testdata/simple_func.ts", root)
	if err != nil {
		t.Fatalf("ParseTypeScriptFile failed: %v", err)
	}

	if len(tf.Funcs) != 1 || tf.Funcs[0] != "greet" {
		t.Errorf("Funcs = %v, want [greet]", tf.Funcs)
	}

	if len(tf.Classes) != 0 {
		t.Errorf("Classes = %v, want []", tf.Classes)
	}

	if len(tf.Vars) != 1 || tf.Vars[0] != "MAX_SIZE" {
		t.Errorf("Vars = %v, want [MAX_SIZE]", tf.Vars)
	}

	if tf.Lines != 5 {
		t.Errorf("Lines = %d, want 5", tf.Lines)
	}

	if tf.MaxDepth != 0 {
		t.Errorf("MaxDepth = %d, want 0", tf.MaxDepth)
	}

	if tf.Control != "sequence" {
		t.Errorf("Control = %q, want %q", tf.Control, "sequence")
	}

	if tf.BodyHash == "" {
		t.Error("BodyHash is empty, want non-empty")
	}
}
