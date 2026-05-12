//ff:func feature=chain type=util control=sequence
//ff:what test: TestHasFeature
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestHasFeature(t *testing.T) {
	var sf model.SourceFile = &model.GoFile{Annotation: &model.Annotation{
		Func: map[string]string{"feature": "validate"},
	}}
	if !hasFeature(sf, "validate") {
		t.Error("expected true for feature=validate")
	}
	if hasFeature(sf, "chain") {
		t.Error("expected false for feature=chain")
	}
	if hasFeature(&model.GoFile{}, "validate") {
		t.Error("expected false for nil annotation")
	}
}
