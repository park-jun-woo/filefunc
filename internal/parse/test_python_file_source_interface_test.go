//ff:func feature=parse type=util control=sequence
//ff:what test: TestPythonFileSourceInterface
package parse

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestPythonFileSourceInterface(t *testing.T) {
	pf := &model.PythonFile{
		Path:   "test.py",
		Module: "test",
	}

	var sf model.SourceFile = pf

	if sf.GetLang() != "python" {
		t.Errorf("GetLang() = %q, want %q", sf.GetLang(), "python")
	}

	if sf.GetPath() != "test.py" {
		t.Errorf("GetPath() = %q, want %q", sf.GetPath(), "test.py")
	}

	if sf.GetPackage() != "test" {
		t.Errorf("GetPackage() = %q, want %q", sf.GetPackage(), "test")
	}

	if sf.GetHasInit() {
		t.Error("GetHasInit() = true, want false")
	}
}
