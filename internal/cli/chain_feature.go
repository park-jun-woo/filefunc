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
		prompt, _ := cmd.Flags().GetString("prompt")
		rate, _ := cmd.Flags().GetFloat64("rate")
		scoreModel, _ := cmd.Flags().GetString("model")
		scoreEndpoint, _ := cmd.Flags().GetString("score-endpoint")

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

		funcs := chain.FilterByFeature(files, feature)
		if len(funcs) == 0 {
			return fmt.Errorf("no funcs found for feature=%s", feature)
		}

		metaFlags := chain.ParseMetaFlags(metaRaw)
		var fileMap map[string]*model.GoFile
		if len(metaFlags) > 0 || prompt != "" {
			fileMap = chain.BuildFuncFileMap(files)
		}

		fmt.Fprintf(os.Stdout, "feature=%s (%d funcs)\n\n", feature, len(funcs))
		for _, name := range funcs {
			results := chain.TraverseChon(g, name, chon)

			var scores map[int]float64
			removed := 0
			if prompt != "" {
				scores, err = chain.ScoreRelevance(results, prompt, scoreEndpoint, scoreModel, fileMap)
				if err != nil {
					return err
				}
				effectiveRate := rate
				if effectiveRate == 0 {
					effectiveRate = 0.8
				}
				results, scores, removed = chain.FilterByRate(results, scores, effectiveRate)
			}

			chain.FormatChain(os.Stdout, name, results, metaFlags, fileMap, scores, removed)
			fmt.Fprintln(os.Stdout)
		}
		return nil
	},
}

func init() {
	chainFeatureCmd.Flags().Int("chon", 1, "chon distance (1~3)")
	chainFeatureCmd.Flags().String("meta", "", "annotation metadata to include (meta,what,why,checked,all)")
	chainFeatureCmd.Flags().String("prompt", "", "user task intent for relevance scoring")
	chainFeatureCmd.Flags().Float64("rate", 0, "relevance score threshold (0.0~1.0)")
	chainFeatureCmd.Flags().String("model", "Qwen/Qwen3-Reranker-0.6B", "reranker model name")
	chainFeatureCmd.Flags().String("score-endpoint", "http://localhost:8000", "vLLM endpoint for reranker")
	chainCmd.AddCommand(chainFeatureCmd)
}
