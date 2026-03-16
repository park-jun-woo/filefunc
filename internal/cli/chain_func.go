//ff:func feature=cli type=command control=sequence
//ff:what chain func 서브커맨드 — 특정 func의 호출 관계 추적
package cli

import (
	"os"

	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/spf13/cobra"
)

var chainFuncCmd = &cobra.Command{
	Use:   "func <name>",
	Short: "Trace call chain for a specific func",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		root, _ := cmd.Flags().GetString("root")
		chon, _ := cmd.Flags().GetInt("chon")
		childDepth, _ := cmd.Flags().GetInt("child-depth")
		parentDepth, _ := cmd.Flags().GetInt("parent-depth")

		if err := CheckProjectRoot(root); err != nil {
			return err
		}

		g, _, err := BuildGraph(root)
		if err != nil {
			return err
		}

		var results []chain.ChonResult
		if childDepth > 0 {
			results = chain.TraverseDepth(g, target, "calls", childDepth)
		} else if parentDepth > 0 {
			results = chain.TraverseDepth(g, target, "called-by", parentDepth)
		} else {
			results = chain.TraverseChon(g, target, chon)
		}

		chain.FormatChain(os.Stdout, target, results)
		return nil
	},
}

func init() {
	chainFuncCmd.Flags().Int("chon", 1, "chon distance (1~3)")
	chainFuncCmd.Flags().Int("child-depth", 0, "trace children only to this depth")
	chainFuncCmd.Flags().Int("parent-depth", 0, "trace parents only to this depth")
	chainCmd.AddCommand(chainFuncCmd)
}
