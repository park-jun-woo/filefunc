//ff:func feature=cli type=command control=sequence
//ff:what chain 서브커맨드 부모 정의
package cli

import "github.com/spf13/cobra"

var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Trace func call chains and relationships",
}

func init() {
	chainCmd.PersistentFlags().String("root", ".", "project root (must contain go.mod and codebook.yaml)")
	rootCmd.AddCommand(chainCmd)
}
