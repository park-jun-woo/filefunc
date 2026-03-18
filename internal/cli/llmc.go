//ff:func feature=cli type=command control=sequence
//ff:what llmc 서브커맨드 정의 — LLM으로 what-body 일치 검증 및 checked 서명
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/filefunc/internal/llm"
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
	"github.com/park-jun-woo/filefunc/internal/walk"
	"github.com/spf13/cobra"
)

var llmcCmd = &cobra.Command{
	Use:   "llmc [project-root]",
	Short: "Verify //ff:what matches func body using LLM",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "."
		if len(args) > 0 {
			root = args[0]
		}

		if err := CheckProjectRoot(root); err != nil {
			return err
		}

		providerName, _ := cmd.Flags().GetString("provider")
		modelName, _ := cmd.Flags().GetString("model")
		endpoint, _ := cmd.Flags().GetString("endpoint")
		threshold, _ := cmd.Flags().GetFloat64("threshold")

		if err := llm.CheckModel(endpoint, modelName); err != nil {
			return err
		}

		provider, err := llm.NewProvider(providerName, endpoint, modelName)
		if err != nil {
			return err
		}

		ignorePatterns := walk.ParseFFIgnore(filepath.Join(root, ".ffignore"))
		paths, err := walk.WalkGoFiles(root, ignorePatterns)
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

		passed, failed, skipped := 0, 0, 0
		for _, gf := range files {
			if gf.IsTest || len(gf.Funcs) == 0 || gf.Annotation == nil {
				continue
			}
			switch ProcessLlmcFile(gf, provider, modelName, threshold) {
			case "skip":
				skipped++
			case "pass":
				passed++
			case "fail":
				failed++
			}
		}

		fmt.Printf("\nllmc: %d passed, %d failed, %d skipped\n", passed, failed, skipped)
		if failed > 0 {
			return fmt.Errorf("%d file(s) failed verification", failed)
		}
		return nil
	},
}

func init() {
	llmcCmd.Flags().String("provider", EnvOrDefault("FILEFUNC_LLM_PROVIDER", "ollama"), "LLM provider")
	llmcCmd.Flags().String("model", EnvOrDefault("FILEFUNC_LLM_MODEL", "gpt-oss:20b"), "LLM model name")
	llmcCmd.Flags().String("endpoint", EnvOrDefault("FILEFUNC_LLM_ENDPOINT", "http://localhost:11434"), "LLM API endpoint")
	llmcCmd.Flags().Float64("threshold", 0.8, "minimum score for passing (0.0~1.0)")
	rootCmd.AddCommand(llmcCmd)
}
