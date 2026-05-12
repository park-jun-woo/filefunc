//ff:func feature=parse type=util control=sequence
//ff:what test: TestParseTypeScriptFileInterface
package parse

import (
	"os/exec"
	"testing"
)

func TestParseTypeScriptFileInterface(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not found, skipping TypeScript parse test")
	}

	root := findTsProjectRoot()
	tf, err := ParseTypeScriptFile("testdata/with_interface.ts", root)
	if err != nil {
		t.Fatalf("ParseTypeScriptFile failed: %v", err)
	}

	if len(tf.Interfaces) != 1 || tf.Interfaces[0] != "Config" {
		t.Errorf("Interfaces = %v, want [Config]", tf.Interfaces)
	}

	if len(tf.TypeAliases) != 1 || tf.TypeAliases[0] != "Options" {
		t.Errorf("TypeAliases = %v, want [Options]", tf.TypeAliases)
	}

	if len(tf.Funcs) != 1 || tf.Funcs[0] != "createConfig" {
		t.Errorf("Funcs = %v, want [createConfig]", tf.Funcs)
	}

	types := tf.GetTypes()
	if len(types) != 2 {
		t.Errorf("GetTypes() len = %d, want 2 (1 interface + 1 type alias)", len(types))
	}
}
