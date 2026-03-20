//ff:func feature=context type=util control=sequence
//ff:what test: TestFilterFeature
package context

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestFilterFeature(t *testing.T) {
	files := []*model.GoFile{
		{Annotation: &model.Annotation{Func: map[string]string{"feature": "validate"}}},
		{Annotation: &model.Annotation{Func: map[string]string{"feature": "chain"}}},
		{Annotation: nil},
	}
	got := FilterFeature(files, []string{"validate"})
	if len(got) != 1 {
		t.Errorf("FilterFeature len = %d, want 1", len(got))
	}
}
