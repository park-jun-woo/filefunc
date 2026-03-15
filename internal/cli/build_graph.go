//ff:func feature=cli type=util
//ff:what 프로젝트 루트에서 호출 그래프를 구성
package cli

import (
	"fmt"

	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
	"github.com/park-jun-woo/filefunc/internal/walk"
)

// BuildGraph builds a call graph from the project at root.
func BuildGraph(root string) (*chain.CallGraph, []*model.GoFile, error) {
	goModPath := FindGoMod(root)
	modulePath, err := parse.ReadModulePath(goModPath)
	if err != nil {
		return nil, nil, fmt.Errorf("reading go.mod: %w", err)
	}

	paths, err := walk.WalkGoFiles(root)
	if err != nil {
		return nil, nil, fmt.Errorf("walking files: %w", err)
	}

	var files []*model.GoFile
	for _, p := range paths {
		gf, err := parse.ParseGoFile(p)
		if err != nil {
			continue
		}
		files = append(files, gf)
	}

	projFuncs, _ := parse.CollectProjectSymbols(files)
	return chain.BuildCallGraph(files, modulePath, projFuncs), files, nil
}
