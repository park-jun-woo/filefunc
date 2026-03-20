//ff:func feature=validate type=util control=sequence
//ff:what test: mustParse
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

func mustParse(t *testing.T, path string) *model.GoFile {
	t.Helper()
	gf, err := parse.ParseGoFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return gf
}
