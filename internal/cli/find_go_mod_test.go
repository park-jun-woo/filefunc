//ff:func feature=cli type=util control=sequence
//ff:what test: TestFindGoMod
package cli

import (
	"os"
	"testing"
)

func TestFindGoMod(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(tmp+"/go.mod", []byte("module test"), 0644)
	os.MkdirAll(tmp+"/sub/deep", 0755)

	got := FindGoMod(tmp + "/sub/deep")
	if got != tmp+"/go.mod" {
		t.Errorf("FindGoMod = %q, want %q", got, tmp+"/go.mod")
	}
}
