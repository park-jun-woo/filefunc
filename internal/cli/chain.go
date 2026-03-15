//ff:func feature=cli type=command
//ff:what chain 서브커맨드 부모 정의
package cli

import "github.com/spf13/cobra"

var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Trace func call chains and relationships",
}

func init() {
	rootCmd.AddCommand(chainCmd)
}
