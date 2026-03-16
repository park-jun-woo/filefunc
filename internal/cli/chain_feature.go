//ff:func feature=cli type=command control=sequence
//ff:what chain feature 서브커맨드 — feature 전체 func의 호출 관계 추적
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/spf13/cobra"
)

var chainFeatureCmd = &cobra.Command{
	Use:   "feature <name>",
	Short: "Trace call chains for all funcs in a feature",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		feature := args[0]
		root, _ := cmd.Flags().GetString("root")
		chon, _ := cmd.Flags().GetInt("chon")
		metaRaw, _ := cmd.Flags().GetString("meta")

		if err := CheckProjectRoot(root); err != nil {
			return err
		}

		g, files, err := BuildGraph(root)
		if err != nil {
			return err
		}

		funcs := chain.FilterByFeature(files, feature)
		if len(funcs) == 0 {
			return fmt.Errorf("no funcs found for feature=%s", feature)
		}

		metaFlags := chain.ParseMetaFlags(metaRaw)
		var fileMap map[string]*model.GoFile
		if len(metaFlags) > 0 {
			fileMap = chain.BuildFuncFileMap(files)
		}

		fmt.Fprintf(os.Stdout, "feature=%s (%d funcs)\n\n", feature, len(funcs))
		for _, name := range funcs {
			results := chain.TraverseChon(g, name, chon)
			chain.FormatChain(os.Stdout, name, results, metaFlags, fileMap)
			fmt.Fprintln(os.Stdout)
		}
		return nil
	},
}

func init() {
	chainFeatureCmd.Flags().Int("chon", 1, "chon distance (1~3)")
	chainFeatureCmd.Flags().String("meta", "", "annotation metadata to include (meta,what,why,checked,all)")
	chainCmd.AddCommand(chainFeatureCmd)
}
