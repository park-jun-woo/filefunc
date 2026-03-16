//ff:func feature=cli type=command control=sequence
//ff:what cobra rootCmd 정의 및 CLI 실행 엔트리포인트
//ff:checked llm=gpt-oss:20b hash=3875c523
package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "filefunc",
	Short: "Code structure convention and toolchain for LLM-native development",
}

func Execute() error {
	return rootCmd.Execute()
}
