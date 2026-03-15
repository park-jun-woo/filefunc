//ff:func feature=cli type=command
//ff:what annotate 서브커맨드 정의 — calls/uses 자동 산출
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/annotate"
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
	"github.com/park-jun-woo/filefunc/internal/walk"
	"github.com/spf13/cobra"
)

var annotateCmd = &cobra.Command{
	Use:   "annotate [path]",
	Short: "Auto-generate //ff:calls and //ff:uses annotations",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "."
		if len(args) > 0 {
			target = args[0]
		}

		goModPath := FindGoMod(target)
		modulePath, err := parse.ReadModulePath(goModPath)
		if err != nil {
			return fmt.Errorf("reading go.mod: %w", err)
		}

		paths, err := walk.WalkGoFiles(target)
		if err != nil {
			return fmt.Errorf("walking files: %w", err)
		}

		var files []*model.GoFile
		for _, p := range paths {
			gf, err := parse.ParseGoFile(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", p, err)
				continue
			}
			files = append(files, gf)
		}

		projFuncs, projTypes := parse.CollectProjectSymbols(files)
		modified := 0

		for _, gf := range files {
			if gf.IsTest || len(gf.Funcs) == 0 {
				continue
			}
			calls, _ := parse.ExtractCalls(gf.Path, modulePath, projFuncs)
			uses, _ := parse.ExtractUses(gf.Path, modulePath, projTypes)

			c1, _ := annotate.WriteAnnotationLine(gf.Path, "calls", strings.Join(calls, ", "))
			c2, _ := annotate.WriteAnnotationLine(gf.Path, "uses", strings.Join(uses, ", "))
			if c1 || c2 {
				modified++
			}
		}

		fmt.Printf("%d file(s) updated.\n", modified)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(annotateCmd)
}
