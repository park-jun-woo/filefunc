//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what PythonFile 슬라이스에서 모듈 레벨 import 방향 그래프를 구축
package validate

import (
	"path/filepath"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// buildImportGraph builds a directed graph of module-level imports among project files.
// Key: absolute file path, Value: list of absolute file paths that the key imports.
// External packages (unresolved imports) are excluded.
func buildImportGraph(files []*model.PythonFile, root string) map[string][]string {
	graph := make(map[string][]string)
	for _, pf := range files {
		absPath, err := filepath.Abs(pf.Path)
		if err != nil {
			continue
		}
		graph[absPath] = resolveFileImports(absPath, pf.ModuleImports, root)
	}
	return graph
}
