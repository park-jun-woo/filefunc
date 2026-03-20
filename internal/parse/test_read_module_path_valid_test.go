//ff:func feature=parse type=util control=sequence
//ff:what test: TestReadModulePath_Valid
package parse

import (
	"testing"
)

func TestReadModulePath_Valid(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/go.mod"
	if err := writeTestFile(path, "module github.com/test/proj\n\ngo 1.22\n"); err != nil {
		t.Fatal(err)
	}
	mod, err := ReadModulePath(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mod != "github.com/test/proj" {
		t.Errorf("module = %q, want %q", mod, "github.com/test/proj")
	}
}
