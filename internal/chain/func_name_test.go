//ff:func feature=chain type=util control=iteration dimension=1
//ff:what test: TestFuncName
package chain

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestFuncName(t *testing.T) {
	tests := []struct {
		name string
		gf   *model.GoFile
		want string
	}{
		{"func", &model.GoFile{Funcs: []string{"Foo"}}, "Foo"},
		{"method", &model.GoFile{Methods: []string{"Bar"}}, "Bar"},
		{"func over method", &model.GoFile{Funcs: []string{"Foo"}, Methods: []string{"Bar"}}, "Foo"},
		{"empty", &model.GoFile{}, ""},
	}
	for _, tt := range tests {
		got := funcName(tt.gf)
		if got != tt.want {
			t.Errorf("funcName(%s) = %q, want %q", tt.name, got, tt.want)
		}
	}
}
