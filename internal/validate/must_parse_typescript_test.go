//ff:func feature=validate type=util control=sequence
//ff:what test: mustParseTypeScript — TypeScript 파일을 파싱하여 TypeScriptFile 반환하는 테스트 헬퍼
package validate

import (
	"os/exec"
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

func mustParseTypeScript(t *testing.T, path string) *model.TypeScriptFile {
	t.Helper()
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not found, skipping TypeScript test")
	}
	root := findTsProjectRoot()
	tf, err := parse.ParseTypeScriptFile(path, root)
	if err != nil {
		t.Fatal(err)
	}
	ann, err := parse.ParseTypeScriptAnnotation(path)
	if err != nil {
		t.Fatal(err)
	}
	tf.Annotation = ann
	return tf
}
