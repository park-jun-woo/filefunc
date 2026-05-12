//ff:func feature=parse type=util control=sequence
//ff:what test: TestTypeScriptFileSourceInterface
package parse

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestTypeScriptFileSourceInterface(t *testing.T) {
	tf := &model.TypeScriptFile{
		Path:   "test.ts",
		Module: "test",
	}

	var sf model.SourceFile = tf

	if sf.GetLang() != "typescript" {
		t.Errorf("GetLang() = %q, want %q", sf.GetLang(), "typescript")
	}

	if sf.GetPath() != "test.ts" {
		t.Errorf("GetPath() = %q, want %q", sf.GetPath(), "test.ts")
	}

	if sf.GetPackage() != "test" {
		t.Errorf("GetPackage() = %q, want %q", sf.GetPackage(), "test")
	}

	if sf.GetHasInit() {
		t.Error("GetHasInit() = true, want false")
	}
}
