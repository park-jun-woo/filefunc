//ff:func feature=chain type=util control=sequence
//ff:what test: TestHasFeature
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestHasFeature(t *testing.T) {
	gf := &model.GoFile{Annotation: &model.Annotation{
		Func: map[string]string{"feature": "validate"},
	}}
	if !hasFeature(gf, "validate") {
		t.Error("expected true for feature=validate")
	}
	if hasFeature(gf, "chain") {
		t.Error("expected false for feature=chain")
	}
	if hasFeature(&model.GoFile{}, "validate") {
		t.Error("expected false for nil annotation")
	}
}
