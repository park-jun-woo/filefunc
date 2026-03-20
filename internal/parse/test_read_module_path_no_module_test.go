//ff:func feature=parse type=util control=sequence
//ff:what test: TestReadModulePath_NoModule
package parse

import (
	"testing"
)

func TestReadModulePath_NoModule(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/go.mod"
	if err := writeTestFile(path, "go 1.22\n"); err != nil {
		t.Fatal(err)
	}
	_, err := ReadModulePath(path)
	if err == nil {
		t.Error("expected error for go.mod without module directive")
	}
}
