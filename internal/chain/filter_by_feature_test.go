//ff:func feature=chain type=util control=sequence
//ff:what test: TestFilterByFeature
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestFilterByFeature(t *testing.T) {
	files := []*model.GoFile{
		{Package: "v", Funcs: []string{"RuleF1"}, Annotation: &model.Annotation{Func: map[string]string{"feature": "validate"}}},
		{Package: "c", Funcs: []string{"Build"}, Annotation: &model.Annotation{Func: map[string]string{"feature": "chain"}}},
	}
	got := FilterByFeature(files, "validate")
	if len(got) != 1 || got[0] != "v.RuleF1" {
		t.Errorf("FilterByFeature = %v, want [v.RuleF1]", got)
	}
}
