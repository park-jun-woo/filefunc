//ff:func feature=validate type=util control=sequence
//ff:what test: mustParsePython
package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

func mustParsePython(t *testing.T, path string) *model.PythonFile {
	t.Helper()
	pf, err := parse.ParsePythonFile(path)
	if err != nil {
		t.Fatal(err)
	}
	ann, err := parse.ParsePythonAnnotation(path)
	if err != nil {
		t.Fatal(err)
	}
	pf.Annotation = ann
	return pf
}
