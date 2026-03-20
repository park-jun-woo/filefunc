//ff:func feature=context type=util control=sequence
//ff:what test: TestMatchAnnotation
package context

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestMatchAnnotation(t *testing.T) {
	ann := &model.Annotation{
		Func: map[string]string{"feature": "validate", "type": "rule"},
	}
	if !matchAnnotation(ann, map[string]string{"feature": "validate"}) {
		t.Error("expected match")
	}
	if matchAnnotation(ann, map[string]string{"feature": "chain"}) {
		t.Error("expected no match")
	}
	if !matchAnnotation(ann, map[string]string{"feature": "validate", "type": "rule"}) {
		t.Error("expected multi-key match")
	}
}
