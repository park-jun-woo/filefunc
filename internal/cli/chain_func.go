//ff:func feature=cli type=command control=sequence
//ff:what chain func 서브커맨드 — 특정 func의 호출 관계 추적
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/park-jun-woo/filefunc/internal/model"
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
		metaRaw, _ := cmd.Flags().GetString("meta")
		prompt, _ := cmd.Flags().GetString("prompt")
		rate, _ := cmd.Flags().GetFloat64("rate")
		scoreModel, _ := cmd.Flags().GetString("model")
		scoreEndpoint, _ := cmd.Flags().GetString("score-endpoint")
		pkg, _ := cmd.Flags().GetString("package")

		if rate > 0 && prompt == "" {
			return fmt.Errorf("--prompt is required when --rate is specified")
		}

		if err := CheckProjectRoot(root); err != nil {
			return err
		}

		g, files, err := BuildGraph(root)
		if err != nil {
			return err
		}

		qualified, err := resolveTarget(g, target, pkg)
		if err != nil {
			return err
		}

		var results []chain.ChonResult
		if childDepth > 0 {
			results = chain.TraverseDepth(g, qualified, "calls", childDepth)
		} else if parentDepth > 0 {
			results = chain.TraverseDepth(g, qualified, "called-by", parentDepth)
		} else {
			results = chain.TraverseChon(g, qualified, chon)
		}

		if pkg != "" {
			results = chain.FilterByPackage(results, pkg)
		}

		metaFlags := chain.ParseMetaFlags(metaRaw)
		var fileMap map[string]*model.GoFile
		if len(metaFlags) > 0 || prompt != "" {
			fileMap = chain.BuildFuncFileMap(files)
		}

		var scores map[int]float64
		removed := 0
		if prompt != "" {
			scores, err = chain.ScoreRelevance(results, prompt, scoreEndpoint, scoreModel, fileMap)
			if err != nil {
				return err
			}
			if rate == 0 {
				rate = 0.8
			}
			results, scores, removed = chain.FilterByRate(results, scores, rate)
		}

		chain.FormatChain(os.Stdout, qualified, results, metaFlags, fileMap, scores, removed)
		return nil
	},
}

func init() {
	chainFuncCmd.Flags().Int("chon", 1, "chon distance (1~3)")
	chainFuncCmd.Flags().Int("child-depth", 0, "trace children only to this depth")
	chainFuncCmd.Flags().Int("parent-depth", 0, "trace parents only to this depth")
	chainFuncCmd.Flags().String("meta", "", "annotation metadata to include (meta,what,why,checked,all)")
	chainFuncCmd.Flags().String("prompt", "", "user task intent for relevance scoring")
	chainFuncCmd.Flags().Float64("rate", 0, "relevance score threshold (0.0~1.0)")
	chainFuncCmd.Flags().String("model", "Qwen/Qwen3-Reranker-0.6B", "reranker model name")
	chainFuncCmd.Flags().String("score-endpoint", "http://localhost:8000", "vLLM endpoint for reranker")
	chainFuncCmd.Flags().String("package", "", "limit to funcs in this Go package")
	chainCmd.AddCommand(chainFuncCmd)
}
