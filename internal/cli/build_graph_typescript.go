//ff:func feature=cli type=util control=sequence
//ff:what TypeScript 프로젝트 루트에서 호출 그래프를 구성
package cli

import (
	"path/filepath"

	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
	"github.com/park-jun-woo/filefunc/internal/walk"
)

// BuildGraphTypeScript builds a call graph from a TypeScript project at root.
func BuildGraphTypeScript(root string) (*chain.CallGraph, []model.SourceFile, error) {
	ignorePatterns := walk.ParseFFIgnore(filepath.Join(root, ".ffignore"))
	files, err := LoadTypeScriptFiles(root, ignorePatterns)
	if err != nil {
		return nil, nil, err
	}

	projFuncs, _ := parse.CollectProjectSymbols(files)
	g := chain.BuildCallGraph(files, "", projFuncs)
	return g, files, nil
}
