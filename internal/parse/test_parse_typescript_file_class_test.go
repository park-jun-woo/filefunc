//ff:func feature=parse type=util control=iteration dimension=1
//ff:what test: TestParseTypeScriptFileClass
package parse

import (
	"os/exec"
	"testing"
)

func TestParseTypeScriptFileClass(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not found, skipping TypeScript parse test")
	}

	root := findTsProjectRoot()
	tf, err := ParseTypeScriptFile("testdata/with_class.ts", root)
	if err != nil {
		t.Fatalf("ParseTypeScriptFile failed: %v", err)
	}

	if len(tf.Classes) != 1 || tf.Classes[0] != "Server" {
		t.Errorf("Classes = %v, want [Server]", tf.Classes)
	}

	if len(tf.Methods) != 2 {
		t.Errorf("Methods count = %d, want 2", len(tf.Methods))
	}

	if !tf.HasConstructor {
		t.Error("HasConstructor = false, want true")
	}

	if len(tf.Funcs) != 0 {
		t.Errorf("Funcs = %v, want []", tf.Funcs)
	}

	wantMethods := map[string]bool{"Server.start": true, "Server.stop": true}
	for _, m := range tf.Methods {
		if !wantMethods[m] {
			t.Errorf("unexpected method %q", m)
		}
	}
}
